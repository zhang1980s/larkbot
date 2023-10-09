package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"msg-event/model"
	"msg-event/model/event"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	supporttypes "github.com/aws/aws-sdk-go-v2/service/support/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	STATUS_NEW      = "NEW"
	STATUS_OPEN     = "OPEN"
	STATUS_CLOSE    = "CLOSE"
	TYPE_OPEN_CASE  = "OPEN_CASE"
	TYPE_CASE       = "CASE"
	SK              = "AWS_CASE"
	GSI_NAME        = "status-type-index"
	GSI_CREATE_TIME = "create-time-index"
	GSI_MSG_ID      = "card_msg_id-index"
)

var tableName = os.Getenv("CASES_TABLE")
var DBClient *dynamodb.Client

func GetDBClient() *dynamodb.Client {
	if DBClient != nil {
		return DBClient
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	///
	// Set the AWS Region that the service clients should use
	// cfg.Region = os.Getenv("AWS_REGION")

	// Using the Config value, create the DynamoDB client
	DBClient := dynamodb.NewFromConfig(cfg)
	return DBClient
}

// OpenCase every time rewrite the one case from this channel
func OpenCase(fromChannelID, customerID, title, msgID string, msg *model.FeiShuMsg) (c *Case, err error) {

	// insert the data into dynamodb
	ca, err := UpsertCase(&Case{
		UserID:        customerID,
		SortKey:       SK,
		ChannelID:     fromChannelID,
		FromChannelID: fromChannelID,
		CreateTime:    time.Now().String(),
		UpdateTime:    time.Now().String(),
		Title:         title,
		Status:        STATUS_NEW,
		Type:          TYPE_OPEN_CASE,
		CardRespMsgID: msgID,
		CardMsg:       msg,
	})
	if err != nil {
		logrus.Errorf("failed to update case for DDB %+v", err)
		return nil, err
	}
	return ca, nil
}

func UpsertCase(c *Case) (ca *Case, err error) {
	client := GetDBClient()
	item, err := attributevalue.MarshalMap(c)

	if err != nil {
		logrus.Errorf("Marshamap failed %v", err)
	}

	logrus.Infof("item %s", item)
	input := &dynamodb.PutItemInput{
		Item:                   item,
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		TableName:              aws.String(tableName),
	}
	_, err = client.PutItem(context.Background(), input)

	if err != nil {
		logrus.Errorf("failed to put data %v", err)
		return nil, err
	}
	return c, nil
}

func convert(attr map[string]types.AttributeValue) *Case {
	c := &Case{}
	attributevalue.UnmarshalMap(attr, c)
	return c
}

func GetCaseByEvent(e *event.Msg) (c *Case, err error) {
	if e.Event.Message.ChatID != "" {
		return GetCase(e.Event.Message.ChatID)
	}
	if e.OpenMsgID != "" {
		return GetCaseByCardMSGID(e.OpenMsgID)
	}
	return nil, nil
}

func GetCase(channelID string) (c *Case, err error) {
	client := GetDBClient()
	// check existing case
	result, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: channelID},
			"sk": &types.AttributeValueMemberS{Value: SK},
		},
		TableName: aws.String(tableName),
	})

	if err != nil {
		return nil, err
	}
	if result != nil && result.Item != nil {
		return convert(result.Item), nil
	}
	return nil, errors.New("您还没有开case")
}

func GetCaseByCardMSGID(msgID string) (c *Case, err error) {
	client := GetDBClient()

	resp, err := client.Query(context.Background(), &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("#v_card_msg_id = :v1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: msgID},
		},
		ExpressionAttributeNames: map[string]string{
			"#v_card_msg_id": "card_msg_id",
		},
		IndexName: aws.String(GSI_MSG_ID),
		TableName: aws.String(tableName),
	})

	if err != nil {
		logrus.Errorf("failed to list all cases %s", err)
		return nil, err
	}
	if len(resp.Items) >= 1 {
		return convert(resp.Items[0]), nil
	} else {
		logrus.Errorf("msg ID %v, resp %v", msgID, resp)
		return nil, errors.New("没有找到工单卡片")
	}
}

func GetCasesByTime(t string) (cs []*Case, err error) {
	logrus.Infof("Start to get all cases by time: %v", t)
	client := GetDBClient()

	if t == "" {
		t = "7"
	}

	intTime, err := strconv.Atoi(t)
	if err != nil {
		logrus.Errorf("failed to convert time to int %s", err)
		return nil, err
	}

	// limit := 10
	cs = []*Case{}
	params := &dynamodb.ScanInput{
		TableName:        aws.String(tableName),
		FilterExpression: aws.String("create_time >= :start_time"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":start_time": &types.AttributeValueMemberS{
				Value: time.Now().AddDate(0, 0, -intTime).Format("2006-01-02"),
			},
		},
	}
	result := &dynamodb.ScanOutput{}
	for {
		result, err = client.Scan(
			context.Background(),
			params,
		)
		if err != nil {
			logrus.Errorf("failed to list all cases %s", err)
			return nil, err
		}
		items := []*Case{}
		err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
		if err != nil {
			fmt.Println("Error unmarshaling DynamoDB items:", err)
			return nil, err
		}
		// logrus.Infof("Get %v cases completed", limit)
		cs = append(cs, items...)
		if result.LastEvaluatedKey == nil {
			break
		}

		params.ExclusiveStartKey = result.LastEvaluatedKey
	}
	return cs, nil
}

func GetProcessingCases() (cs []*Case, err error) {
	logrus.Infof("Start to get all un-closed cases")
	client := GetDBClient()
	resp, err := client.Query(context.Background(), &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("#v_status = :v1 AND #v_type = :v2"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: STATUS_OPEN},
			":v2": &types.AttributeValueMemberS{Value: TYPE_CASE},
		},
		ExpressionAttributeNames: map[string]string{
			"#v_status": "status",
			"#v_type":   "type",
		},
		IndexName: aws.String(GSI_NAME),
		TableName: aws.String(tableName),
	})
	if err != nil {
		logrus.Errorf("failed to list all cases %s", err)
		return nil, err
	}
	cs = make([]*Case, len(resp.Items))
	for i, v := range resp.Items {
		cs[i] = convert(v)
	}
	logrus.Infof("Get all un-closed cases completed")
	return cs, nil
}

type Case struct {
	AccountKey      string    `dynamodbav:"account_key"`
	UserID          string    `dynamodbav:"user_id"`
	ChannelID       string    `dynamodbav:"pk"`
	SortKey         string    `dynamodbav:"sk"`
	FromChannelID   string    `dynamodbav:"from_channel"`
	CreateTime      string    `dynamodbav:"create_time"`
	UpdateTime      string    `dynamodbav:"update_time"`
	Title           string    `dynamodbav:"title"`
	CaseID          string    `dynamodbav:"case_id"`
	CaseURL         string    `dynamodbav:"case_url"`
	Content         string    `dynamodbav:"content"`
	Status          string    `dynamodbav:"status"`
	ServiceCode     string    `dynamodbav:"service_code"`
	SevCode         string    `dynamodbav:"sev_code"`
	Type            string    `dynamodbav:"type"`
	LastCommentTime time.Time `dynamodbav:"last_comment_time"`
	Comments        []supporttypes.Communication
	DisplayCaseID   string           `dynamodbav:"display_case_id"`
	CardRespMsgID   string           `dynamodbav:"card_msg_id"`
	CardMsg         *model.FeiShuMsg `dynamodbav:"card_msg"`
}

// GetKey returns the primary key of the case in a format that can be
// sent to DynamoDB.
func (c Case) GetKey() map[string]types.AttributeValue {
	pk, err := attributevalue.Marshal(c.ChannelID)
	if err != nil {
		logrus.Errorf("failed to get case key when convert : %v", err)
	}
	sk, err := attributevalue.Marshal(c.SortKey)
	if err != nil {
		logrus.Errorf("failed to get case key when convert : %v", err)
	}
	return map[string]types.AttributeValue{"pk": pk, "sk": sk}
}

func (c Case) Print() {
	str, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		logrus.Errorf("print failed %v", err)
	}
	logrus.Infoln(string(str))
}
