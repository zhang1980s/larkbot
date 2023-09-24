package dao

import (
	"context"
	"errors"
	"lambda-lark-bot/config"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

var cfgTableName = "bot_config"
var EnvConfigKey = "CFG_KEY"
var AccessKeyID = "XXXX"
var SecretAccessKey = "YYYY"
var AppID = "cli_xxxx"
var AppSecret = "yyyy"

func SetupConfig() error {
	client := GetDBClient()
	// check existing request
	result, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: os.Getenv(EnvConfigKey)},
		},
		TableName: aws.String(cfgTableName),
	})
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
		item, err := attributevalue.MarshalMap(c)

		if err != nil {
			logrus.Errorf("Marshal map failed %v", err)
		}
		logrus.Infof("item is %s", item)
		input := &dynamodb.PutItemInput{
			Item:                   item,
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
			TableName:              aws.String(cfgTableName),
		}
		_, err = client.PutItem(context.Background(), input)

		if err != nil {
			logrus.Errorf("failed to put data %v", err)
			return errors.New(config.BotConfigNotExisted)
		}
	}
	return nil
}

func convertCfg(attr map[string]types.AttributeValue) *config.Config {
	c := &config.Config{}

	err := attributevalue.UnmarshalMap(attr, c)

	if err != nil {
		logrus.Errorf("failed to unmarshal map %v", err)
	}
	return c
}
