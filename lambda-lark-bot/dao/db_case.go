package dao

import (
	"errors"
	"lambda-lark-bot/model"
	"lambda-lark-bot/model/event"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go-v2/service/support"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	tableName       = "bot_cases"
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

var DBClient *dynamodb.Client

func GetDBClient() *dynamodb.Client {
	if DBClient != nil {
		return DBClient
	}
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	// Set the AWS Region that the service clients should use
	cfg.Region = os.Getenv("AWS_REGION")

	// Using the Config value, create the DynamoDB client
	DBClient := dynamodb.New(cfg)
	return DBClient
}

//OpenCase every time rewrite the one case from this channel
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
	item, err := dynamodbattribute.MarshalMap(c)

	if err != nil {
		logrus.Errorf("Marshamap failed %v", err)
	}

	logrus.Infof("item %s", item)
	input := &dynamodb.PutItemInput{
		Item:                   item,
		ReturnConsumedCapacity: dynamodb.ReturnConsumedCapacityTotal,
		TableName:              aws.String(tableName),
	}
	req := client.PutItemRequest(input)
	_, err = req.Send(context.Background())
	if err != nil {
		logrus.Errorf("failed to put data %v", err)
		return nil, err
	}
	return c, nil
}

func convert(attr map[string]dynamodb.AttributeValue) *Case {
	c := &Case{}
	dynamodbattribute.UnmarshalMap(attr, c)
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
	req := client.GetItemRequest(&dynamodb.GetItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"pk": {S: aws.String(channelID)},
			"sk": {S: aws.String(SK)},
		},
		TableName: aws.String(tableName),
	})
	result, err := req.Send(context.Background())
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
	query := client.QueryRequest(&dynamodb.QueryInput{
		KeyConditionExpression: aws.String("#v_card_msg_id = :v1"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":v1": {
				S: aws.String(msgID),
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#v_card_msg_id": "card_msg_id",
		},
		IndexName: aws.String(GSI_MSG_ID),
		TableName: aws.String(tableName),
	})
	resp, err := query.Send(context.Background())
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

// func GetCasesByTime(time string) (cs []*Case, err error) {
// 	logrus.Infof("Start to get all cases by time: %v", time)
// 	client := GetDBClient()
// 	// Map<String, AttributeValue> lastKeyEvaluated = null;
// 	// do {
// 	// 	ScanRequest scanRequest = new ScanRequest()
// 	// 		.withTableName(tableName)
// 	// 		.withLimit(10)
// 	// 		.withExclusiveStartKey(lastKeyEvaluated);

// 	// 	ScanResult result = client.scan(scanRequest);
// 	// 	for (Map<String, AttributeValue> item : result.getItems()){
// 	// 		printItem(item);
// 	// 	}
// 	// 	lastKeyEvaluated = result.getLastEvaluatedKey();
// 	// } while (lastKeyEvaluated != null);

// 	intTime, err := strconv.Atoi(time)
// 	if err != nil {
// 		logrus.Errorf("failed to convert time to int %s", err)
// 		return nil, err
// 	}

// 	limit := 10
// 	cs = []*Case{}
// 	do {
// 		//FIXME check sdk code
// 		query := client.QueryRequest(&dynamodb.ScanRequest{
// 			FilterExpression: aws.String("create_time >= :create_time"),
// 			//FIXME check the attr
// 			Limit: limit,
// 			ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
// 				":create_time": {
// 					S: aws.String(time.Now().AddDate(0, -intTime, 0).String()),
// 				},
// 			},
// 			IndexName: aws.String(GSI_CREATE_TIME),
// 			TableName: aws.String(tableName),
// 		})
// 		resp, err := query.Send(context.Background())
// 		if err != nil {
// 			logrus.Errorf("failed to list all cases %s", err)
// 			return nil, err
// 		}
// 		rs = make([]*Case, len(resp.Items))
// 		for i, v := range resp.Items {
// 			rs[i] = convert(v)
// 		}
// 		logrus.Infof("Get %v cases completed", limit)
// 		cs = append(cs, rs)
// 	}
// 	return cs, nil
// }

func GetProcessingCases() (cs []*Case, err error) {
	logrus.Infof("Start to get all un-closed cases")
	client := GetDBClient()
	query := client.QueryRequest(&dynamodb.QueryInput{
		KeyConditionExpression: aws.String("#v_status = :v1 AND #v_type = :v2"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":v1": {
				S: aws.String(STATUS_OPEN),
			},
			":v2": {
				S: aws.String(TYPE_CASE),
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#v_status": "status",
			"#v_type":   "type",
		},
		IndexName: aws.String(GSI_NAME),
		TableName: aws.String(tableName),
	})
	resp, err := query.Send(context.Background())
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
	AccountKey      string    `json:"account_key,omitempty"`
	UserID          string    `json:"user_id,omitempty"`
	ChannelID       string    `json:"pk,omitempty"`
	SortKey         string    `json:"sk,omitempty"`
	FromChannelID   string    `json:"from_channel,omitempty"`
	CreateTime      string    `json:"create_time,omitempty"`
	UpdateTime      string    `json:"update_time,omitempty"`
	Title           string    `json:"title,omitempty"`
	CaseID          string    `json:"case_id,omitempty"`
	Content         string    `json:"content,omitempty"`
	Status          string    `json:"status,omitempty"`
	ServiceCode     string    `json:"service_code,omitempty"`
	SevCode         string    `json:"sev_code,omitempty"`
	Type            string    `json:"type,omitempty"`
	LastCommentTime time.Time `json:"last_comment_time,omitempty"`
	Comments        []support.Communication
	DisplayCaseID   string           `json:"display_case_id,omitempty"`
	CardRespMsgID   string           `json:"card_msg_id,omitempty"`
	CardMsg         *model.FeiShuMsg `json:"card_msg,omitempty"`
}
