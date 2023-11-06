package dao

import (
	"fmt"
	"msg-event/config"
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

	var updateExpParts []string
	attrNames := map[string]string{"#UserWhiteListMap": "user_whitelist"}
	attrValues := map[string]types.AttributeValue{}

	for key, value := range whitelist {
		attrKey := "#K_" + key
		attrValue := ":V_" + key

		updateExpParts = append(updateExpParts, fmt.Sprintf("#UserWhiteListMap.%s = %s", attrKey, attrValue))
		attrNames[attrKey] = key
		attrValues[attrValue] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%s", value)}
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
		logrus.Errorf("Failed to update item, %v", err)
	}
	return err
}

func DelWhiteList(whiteList map[string]string) (err error) {
	client := GetDBClient()
	primaryKeyValue := os.Getenv("CFG_KEY")

	var updateExpParts []string
	attrNames := map[string]string{"#UserWhiteListMap": "user_whitelist"}

	for key, _ := range whiteList {
		attrKey := "#K_" + key
		updateExpParts = append(updateExpParts, fmt.Sprintf("#UserWhiteListMap.%s", attrKey))
		attrNames[attrKey] = key
		if _, ok := config.Conf.RoleMap[key]; ok {
			attrNames["#RoleMap"] = "role"
			updateExpParts = append(updateExpParts, fmt.Sprintf("#RoleMap.%s", attrKey))
		}
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                aws.String(cfgTableName),
		Key:                      map[string]types.AttributeValue{"key": &types.AttributeValueMemberS{Value: primaryKeyValue}}, // Replace with your primary key attribute name and value
		UpdateExpression:         aws.String("REMOVE " + strings.Join(updateExpParts, ", ")),
		ExpressionAttributeNames: attrNames,
	}

	_, err = client.UpdateItem(context.TODO(), input)
	if err != nil {
		logrus.Errorf("Failed to update item, %v", err)
	}

	return err
}

func GetWhiteList() (whiteList map[string]string) {

	whiteList = make(map[string]string)
	userList := make(map[string]string)
	userList = config.Conf.UserWhiteListMap
	roleMap := make(map[string]string)
	roleMap = config.Conf.RoleMap

	for key, value := range userList {
		if _, ok := roleMap[key]; !ok {
			whiteList[value] = "Admin"
		} else {
			whiteList[value] = "User"
		}
	}
	return whiteList
}

func SetAdmin(adminList map[string]string) (err error) {

	client := GetDBClient()
	primaryKeyValue := os.Getenv("CFG_KEY")

	var updateExpParts []string
	attrNames := map[string]string{"#RoleMap": "role"}
	attrValues := map[string]types.AttributeValue{}

	for key, value := range adminList {

		attrKey := "#K_" + key
		attrValue := ":V_" + key

		if _, ok := config.Conf.UserWhiteListMap[key]; !ok {
			attrNames["#UserWhiteListMap"] = "user_whitelist"
			updateExpParts = append(updateExpParts, fmt.Sprintf("#UserWhiteListMap.%s = %s", attrKey, attrValue))
		}

		updateExpParts = append(updateExpParts, fmt.Sprintf("#RoleMap.%s = %s", attrKey, attrValue))
		attrNames[attrKey] = key
		attrValues[attrValue] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%s", value)}
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(cfgTableName),
		Key:                       map[string]types.AttributeValue{"key": &types.AttributeValueMemberS{Value: primaryKeyValue}}, // Replace with your primary key attribute name and value
		UpdateExpression:          aws.String("SET " + strings.Join(updateExpParts, ", ")),
		ExpressionAttributeNames:  attrNames,
		ExpressionAttributeValues: attrValues,
	}

	_, err = client.UpdateItem(context.TODO(), input)
	if err != nil {
		logrus.Errorf("Failed to update item, %v", err)
	}

	return err
}
