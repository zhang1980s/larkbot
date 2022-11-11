package dao

import (
	"errors"
	"lambda-lark-bot/config"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

var cfgTableName = "bot_config"
var EnvConfigKey = "CFG_KEY"
var AccessKeyID = "XXXX"
var SecretAccessKey = "YYYY"
var AppID = "cli_xxxx"
var AppSecret = "yyyy"

func SetupConfig() error {
	client := GetDBClient()
	// check existing case
	req := client.GetItemRequest(&dynamodb.GetItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"key": {S: aws.String(os.Getenv(EnvConfigKey))},
		},
		TableName: aws.String(cfgTableName),
	})
	result, err := req.Send(context.Background())
	if err != nil {
		return err
	}
	if result != nil && result.Item != nil {
		config.Conf = convertCfg(result.Item)
		config.Usage = config.Conf.Usage
		config.SevMap = config.Conf.SevMap
		config.ServiceMap = config.Conf.ServiceMap
	} else {
		c := config.Config{
			Key: os.Getenv(EnvConfigKey),
		}
		c.Usage = config.Usage
		c.ServiceMap = config.ServiceMap
		c.SevMap = config.SevMap
		c.Accounts = map[string]*config.Account{
			"0": {
				AccessKeyID:     AccessKeyID,
				SecretAccessKey: SecretAccessKey,
			},
		}
		c.AppID = AppID
		c.AppSecret = AppSecret
		c.CaseCardTemplate = config.CardTemplate
		c.ErrCardTemplate = config.ErrCardTemplate
		item, err := dynamodbattribute.MarshalMap(c)

		if err != nil {
			logrus.Errorf("Marshal map failed %v", err)
		}
		logrus.Infof("item %s", item)
		input := &dynamodb.PutItemInput{
			Item:                   item,
			ReturnConsumedCapacity: dynamodb.ReturnConsumedCapacityTotal,
			TableName:              aws.String(cfgTableName),
		}
		req := client.PutItemRequest(input)
		_, err = req.Send(context.Background())
		if err != nil {
			logrus.Errorf("failed to put data %v", err)
			return errors.New(config.BotConfigNotExisted)
		}
	}
	return nil
}

func convertCfg(attr map[string]dynamodb.AttributeValue) *config.Config {
	c := &config.Config{}
	dynamodbattribute.UnmarshalMap(attr, c)
	return c
}
