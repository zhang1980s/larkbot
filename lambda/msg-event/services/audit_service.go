package services

import (
	"msg-event/dao"
	"msg-event/model/event"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var auditTableName = os.Getenv("AUDIT_TABLE")

func Processable(e *event.Msg) bool {
	logrus.Infof("messageid is %s", e.Event.Message.MsgID)
	if e.Event.Message.MsgID == "" {
		return true
	}

	client := dao.GetDBClient()
	// check existing request
	result, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: e.Event.Message.MsgID},
		},
		TableName: aws.String(auditTableName),
	})

	if err != nil {
		logrus.Errorf("check audit failed %v", err)
		return true
	}
	if result != nil && result.Item != nil {
		logrus.Infof("duplicate msg %v", result)
		return false
	}

	item := map[string]types.AttributeValue{
		"key": &types.AttributeValueMemberS{Value: e.Event.Message.MsgID},
	}
	logrus.Infof("item %s", item)
	input := &dynamodb.PutItemInput{
		Item:                   item,
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		TableName:              aws.String(auditTableName),
	}

	_, err = client.PutItem(context.Background(), input)
	if err != nil {
		logrus.Errorf("failed to put data %v", err)
	}
	return true
}
