package services

import (
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var auditTableName = "audit"

func Processable(e *event.Msg) bool {
	logrus.Infof("EventID is %s", e.Header.EventID)
	if e.Header.EventID == "" {
		return true
	}

	client := dao.GetDBClient()
	// check existing case
	req := client.GetItemRequest(&dynamodb.GetItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"key": {S: aws.String(e.Header.EventID)},
		},
		TableName: aws.String(auditTableName),
	})
	result, err := req.Send(context.Background())

	if err != nil {
		logrus.Errorf("check audit failed %v", err)
		return true
	}
	if result != nil && result.Item != nil {
		logrus.Infof("duplicate msg %v", result)
		return false
	}

	item := map[string]dynamodb.AttributeValue{
		"key": {S: aws.String(e.Header.EventID)},
	}
	logrus.Infof("item %s", item)
	input := &dynamodb.PutItemInput{
		Item:                   item,
		ReturnConsumedCapacity: dynamodb.ReturnConsumedCapacityTotal,
		TableName:              aws.String(auditTableName),
	}
	r := client.PutItemRequest(input)
	_, err = r.Send(context.Background())
	if err != nil {
		logrus.Errorf("failed to put data %v", err)
	}
	return true
}
