package dao

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func AddWhitelist(whitelist map[string]string) (err error) {
	client := GetDBClient()
	primaryKeyValue := os.Getenv("CFG_KEY")
	role := "user"

	var updateExpParts []string
	attrNames := map[string]string{"#UserWhiteListMap": "user_whitelist"}
	attrValues := map[string]types.AttributeValue{}

	for key, value := range whitelist {
		attrKey := "#K_" + key
		attrValue := ":V_" + role

		updateExpParts = append(updateExpParts, fmt.Sprintf("#UserWhiteListMap.%s = %s", attrKey, attrValue))
		attrNames[attrKey] = key
		attrValues[attrValue] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%s-%s", role, value)}
	}

	// Define the UpdateItem input
	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(cfgTableName),
		Key:                       map[string]types.AttributeValue{"key": &types.AttributeValueMemberS{Value: primaryKeyValue}}, // Replace with your primary key attribute name and value
		UpdateExpression:          aws.String("SET " + strings.Join(updateExpParts, ", ")),
		ExpressionAttributeNames:  attrNames,
		ExpressionAttributeValues: attrValues,
	}

	// Update the item in DynamoDB
	_, err = client.UpdateItem(context.TODO(), input)
	if err != nil {
		// log.Fatalf("Failed to update item, %v", err)
		logrus.Error("Failed to update item, %v", err)
	}
	return err
}

func DelWhiteList(whiteList []string) (err error) {
	client := GetDBClient()
	primaryKeyValue := os.Getenv("CFG_KEY")

	var updateExpParts []string
	attrNames := map[string]string{"#UserWhiteListMap": "user_whitelist"}

	for _, key := range whiteList {
		attrKey := "#K_" + key
		updateExpParts = append(updateExpParts, fmt.Sprintf("#UserWhiteListMap.%s", attrKey))
		attrNames[attrKey] = key
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                aws.String("LarkbotAppStack-botconfigCFB80AE2-VPP11H16RTMN"),
		Key:                      map[string]types.AttributeValue{"key": &types.AttributeValueMemberS{Value: primaryKeyValue}}, // Replace with your primary key attribute name and value
		UpdateExpression:         aws.String("REMOVE " + strings.Join(updateExpParts, ", ")),
		ExpressionAttributeNames: attrNames,
	}

	_, err = client.UpdateItem(context.TODO(), input)
	if err != nil {
		logrus.Error("Failed to update item, %v", err)
	}

	return err
}
