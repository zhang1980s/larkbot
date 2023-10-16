package dao

import (
	"fmt"
	"msg-event/config"
	"os"
	"regexp"
	"time"

	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/support"
	"github.com/aws/aws-sdk-go-v2/service/support/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const caseUrl = "https://support.console.aws.amazon.com/support/home#/case/?displayId=%s"

var SupportClient *support.Client

func GetSupportClient(c *Case) *support.Client {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(getRegion()))
	if err != nil {
		logrus.Printf("Couldn't load default configuration. %v\n", err)
		return nil
	}

	a, ok := config.Conf.Accounts[c.AccountKey]
	if !ok {
		panic("failed to get account " + c.AccountKey)
	}

	stsClient := sts.NewFromConfig(cfg)

	provider := stscreds.NewAssumeRoleProvider(stsClient, a.RoleARN)

	if provider == nil {
		logrus.Errorf("counldn't create assume role provider.")
		panic("failed to get provider by arn " + c.AccountKey)
	}

	cfg.Credentials = provider

	return support.NewFromConfig(cfg)
}

// Deprecated
func GetSupportClientByAKSK(c *Case) *support.Client {
	if SupportClient != nil {
		return SupportClient
	}

	a, ok := config.Conf.Accounts[c.AccountKey]
	if !ok {
		panic("failed to get account " + c.AccountKey)
	}
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(getRegion()),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(a.AccessKeyID, a.SecretAccessKey, "")))

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	logrus.Info("Support Endpoint Region:", cfg.Region)

	SupportClient = support.NewFromConfig(cfg)
	return SupportClient
}

func getRegion() string {
	// Get the value of SUPPORT_REGION
	supportRegion := os.Getenv("SUPPORT_REGION")
	// If it is cn, return cn-north-1
	if supportRegion == "cn" {
		return "cn-north-1"
	}
	// Otherwise, return us-east-1 as default
	return "us-east-1"
}

// Create Case and Create Channel
func CreateCaseAndChannel(c *Case) (*Case, error) {
	client := GetSupportClient(c)
	input := &support.CreateCaseInput{}

	switch os.Getenv("CASE_LANGUAGE") {
	case "zh", "ja", "ko":
		input.Language = aws.String(os.Getenv("CASE_LANGUAGE"))
	default:
		input.Language = aws.String("en")
	}

	input.Subject = &c.Title
	v := config.ServiceMap[c.ServiceCode]
	input.ServiceCode = aws.String(v[0])
	input.CategoryCode = aws.String(v[1])

	v1 := config.SevMap[c.SevCode]
	input.SeverityCode = aws.String(v1)
	input.CommunicationBody = &c.Content

	var response *support.CreateCaseOutput

	err := retry.Do(
		func() error {
			var err error
			response, err = client.CreateCase(context.Background(), input)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		logrus.Errorf("failed to create aws case %s", err)
		return nil, err
	}
	c.CaseID = *response.CaseId
	awsCase, err := GetAWSCase(c)
	if err != nil {
		logrus.Errorf("failed to get aws case %s", err)
		return nil, err
	}

	displayCaseID := awsCase.Cases[0].DisplayId

	c.CaseID = *response.CaseId
	c.DisplayCaseID = *displayCaseID
	c.UpdateTime = time.Now().String()

	logrus.Infof("aws case've been created, then create channel. case id %v", *displayCaseID)

	// careate channel
	channelID, err := CreateChannel([]string{c.UserID}, c.DisplayCaseID+"-"+c.Title)
	if err != nil {
		logrus.Errorf("failed to create feishu channel %s", err)
		return nil, err
	}

	fromChannel := c.ChannelID
	c.FromChannelID = fromChannel
	c.ChannelID = channelID
	c.LastCommentTime = time.Now()
	c.Type = TYPE_CASE

	c.CardMsg.ChatId = c.ChannelID
	c.CardMsg.UserId = c.UserID
	rsp, err := SendCardMsg(c.CardMsg, c)
	if err != nil {
		logrus.Errorf("failed to send case message %s", err)
		return nil, err
	}
	c.CardRespMsgID = *rsp.Data.MessageId

	/// Adding ChatTab with CASE URL
	logrus.Info("Adding ChatTab with CASE URL")
	url := fmt.Sprintf(caseUrl, c.DisplayCaseID)

	err = CreateChatTab(c.ChannelID, url)

	if err != nil {
		logrus.Errorf("failed to create feishu chat tab %s", err)
		return nil, err
	}
	c.CaseURL = url

	a, ok := config.Conf.Accounts[c.AccountKey]
	if !ok {
		panic("failed to get account " + c.AccountKey)
	}

	c.CaseAccountID = GetAccountIdFromRoleARN(a.RoleARN)

	c, err = UpsertCase(c)
	if err != nil {
		logrus.Errorf("failed to update case info %s", err)
		return nil, err
	}

	return c, err
}

func AddAttToCase(c *Case, setID, name string) (caze *Case, err error) {
	client := GetSupportClient(c)

	add := &support.AddCommunicationToCaseInput{
		CaseId:            &c.CaseID,
		AttachmentSetId:   &setID,
		CommunicationBody: &name,
	}
	var resp *support.AddCommunicationToCaseOutput
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.AddCommunicationToCase(context.Background(), add)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		logrus.Errorf("failed to add att %s", err)
		return nil, err
	}

	logrus.Infof("%v", resp)
	return c, nil
}

func AddComment(c *Case, comment string) (caze *Case, err error) {
	client := GetSupportClient(c)
	r := regexp.MustCompile(`^@.+\s+`)
	comment = r.ReplaceAllString(comment, "") // replace @user1

	add := &support.AddCommunicationToCaseInput{
		CaseId:            &c.CaseID,
		CommunicationBody: &comment,
	}
	var resp *support.AddCommunicationToCaseOutput
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.AddCommunicationToCase(context.Background(), add)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		logrus.Errorf("failed to add comment %s", err)
		return nil, err
	}

	logrus.Infof("%v", resp)
	return c, nil
}

func GetAWSCase(c *Case) (caze *support.DescribeCasesOutput, err error) {
	client := GetSupportClient(c)

	input := &support.DescribeCasesInput{
		CaseIdList: []string{c.CaseID},
	}

	var resp *support.DescribeCasesOutput
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.DescribeCases(context.Background(), input)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		logrus.Errorf("failed to get case from aws case %v", resp)
		return nil, err
	}

	return resp, nil
}

func GetCaseComments(c *Case, ltime time.Time) (comments []types.Communication, err error) {
	logrus.Infof("Starting to get case %s comments", *aws.String(c.DisplayCaseID))
	client := GetSupportClient(c)

	input := &support.DescribeCommunicationsInput{
		AfterTime: aws.String(FormatTime(ltime)),
		CaseId:    aws.String(c.CaseID),
	}
	var resp *support.DescribeCommunicationsOutput
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.DescribeCommunications(context.Background(), input)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	logrus.Infof("Get case %s comments complete", *aws.String(c.DisplayCaseID))
	return resp.Communications, nil
}

func AddAttachmentToCase(c *Case, name string, data []byte) error {
	client := GetSupportClient(c)
	att := &support.AddAttachmentsToSetInput{
		AttachmentSetId: nil,
		Attachments: []types.Attachment{
			{
				Data:     data,
				FileName: &name,
			},
		},
	}

	var resp *support.AddAttachmentsToSetOutput
	err := retry.Do(
		func() error {
			var err error
			resp, err = client.AddAttachmentsToSet(context.Background(), att)
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		return err
	}
	logrus.Infof("add att %v", resp)
	_, err = AddAttToCase(c, *resp.AttachmentSetId, "附件："+name)
	if err != nil {
		return err
	}
	return nil
}

func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func FormatComments(comments []types.Communication) string {
	s := ""
	for _, c := range comments {
		s += fmt.Sprintf("来自%s的最新回复(%s):\\n %s\\n", *c.SubmittedBy, *c.TimeCreated, *c.Body)
	}

	return s
}

func FormatMsg(caze *Case) string {
	return fmt.Sprintln("工单创建必要内容缺失。请输入帮助关键字获取使用信息")
}

func GetAccountIdFromRoleARN(s string) string {
	arn := s
	re := regexp.MustCompile(`(?:::)(\d+)(?::)`)
	match := re.FindStringSubmatch(arn)

	if len(match) > 1 {
		result := match[1]
		return result
	} else {
		return "0000"
	}
}
