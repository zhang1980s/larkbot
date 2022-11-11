package dao

import (
	"fmt"
	"lambda-lark-bot/config"
	"os"
	"regexp"
	"time"

	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/support"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var SupportClient *support.Client

func GetSupportClient(c *Case) *support.Client {
	if SupportClient != nil {
		return SupportClient
	}
	logrus.Infof("Conf %v", config.Conf)
	logrus.Infof("case %v", c)
	a, ok := config.Conf.Accounts[c.AccountKey]
	if !ok {
		panic("failed to get account " + c.AccountKey)
	}
	cfg, err := external.LoadDefaultAWSConfig(
		external.WithCredentialsValue(aws.Credentials{
			AccessKeyID:     a.AccessKeyID,
			SecretAccessKey: a.SecretAccessKey,
		}),
	)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use

	cfg.Region = endpoints.UsEast1RegionID

	if os.Getenv("SUPPORT_REGION") == "cn" {
		cfg.Region = endpoints.CnNorth1RegionID
	}

	logrus.Info("Support Endpoint Region:", cfg.Region)

	return support.New(cfg)
}

func CreateCase(c *Case) (*Case, error) {
	client := GetSupportClient(c)
	input := &support.CreateCaseInput{}
	input.Subject = &c.Title
	v := config.ServiceMap[c.ServiceCode]
	input.ServiceCode = aws.String(v[0])
	input.CategoryCode = aws.String(v[1])

	v1 := config.SevMap[c.SevCode]
	input.SeverityCode = aws.String(v1)
	input.CommunicationBody = &c.Content

	var response *support.CreateCaseResponse

	err := retry.Do(
		func() error {
			var err error
			response, err = client.CreateCaseRequest(input).Send(context.Background())
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
	c.CardRespMsgID = rsp.Data.MsgID
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
	var resp *support.AddCommunicationToCaseResponse
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.AddCommunicationToCaseRequest(add).Send(context.Background())
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

	logrus.Infof("%s", resp.String())
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
	var resp *support.AddCommunicationToCaseResponse
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.AddCommunicationToCaseRequest(add).Send(context.Background())
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

	logrus.Infof("%s", resp.String())
	return c, nil
}

func GetAWSCase(c *Case) (caze *support.DescribeCasesResponse, err error) {
	client := GetSupportClient(c)

	input := &support.DescribeCasesInput{
		CaseIdList: []string{c.CaseID},
	}

	var resp *support.DescribeCasesResponse
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.DescribeCasesRequest(input).Send(context.Background())
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		logrus.Errorf("failed to get case from aws case %s", resp)
		return nil, err
	}

	return resp, nil
}

func GetCaseComments(c *Case, ltime time.Time) (comments []support.Communication, err error) {
	client := GetSupportClient(c)

	input := &support.DescribeCommunicationsInput{
		AfterTime: aws.String(FormatTime(ltime)),
		CaseId:    aws.String(c.CaseID),
	}
	var resp *support.DescribeCommunicationsResponse
	err = retry.Do(
		func() error {
			var err error
			resp, err = client.DescribeCommunicationsRequest(input).Send(context.Background())
			if err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		return nil, err
	}
	return resp.Communications, nil
}

func AddAttachmentToCase(c *Case, name string, data []byte) error {
	client := GetSupportClient(c)
	att := &support.AddAttachmentsToSetInput{
		AttachmentSetId: nil,
		Attachments: []support.Attachment{
			{
				Data:     data,
				FileName: &name,
			},
		},
	}

	var resp *support.AddAttachmentsToSetResponse
	err := retry.Do(
		func() error {
			var err error
			resp, err = client.AddAttachmentsToSetRequest(att).Send(context.Background())
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

func FormatComments(comments []support.Communication) string {
	s := ""
	for _, c := range comments {
		s += fmt.Sprintf("来自%s的最新回复(%s):\n %s\n", *c.SubmittedBy, *c.TimeCreated, *c.Body)
	}
	return s
}

func FormatMsg(caze *Case) string {
	return fmt.Sprintln("工单创建必要内容缺失。请输入帮助关键字获取使用信息")
}
