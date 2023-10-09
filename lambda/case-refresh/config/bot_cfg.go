package config

import (
	"case-refresh/model"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

var Conf *Config

type Config struct {
	Key              string              `dynamodbav:"key"`
	SevMap           map[string]string   `dynamodbav:"sev_map"`
	ServiceMap       map[string][]string `dynamodbav:"service_map"`
	Usage            string              `dynamodbav:"usage"`
	Accounts         map[string]*Account `dynamodbav:"accounts"`
	AppID            string              `dynamodbav:"app_id"`
	AppSecret        string              `dynamodbav:"app_secret"`
	AppIDARN         string              `dynamodbav:"app_id_arn"`
	AppSecretARN     string              `dynamodbav:"app_secret_arn"`
	ErrCardTemplate  *model.FeiShuMsg    `dynamodbav:"err_card_template"`
	CaseCardTemplate *model.FeiShuMsg    `dynamodbav:"case_card_template"`
	Ack              string              `dynamodbav:"ack"`
	NoPermissionMSG  string              `dynamodbav:"no_permission_msg"`
	UserWhiteListMap map[string]string   `dynamodbav:"user_whitelist"`
}

type Account struct {
	AccessKeyID     string `dynamodbav:"access_key_id"`
	SecretAccessKey string `dynamodbav:"secret_access_key"`
	RoleARN         string `dynamodbav:"role_arn"`
}

// GetKey returns the primary key of the cfg in a format that can be
// sent to DynamoDB.
func (c Config) GetKey() map[string]types.AttributeValue {
	key, err := attributevalue.Marshal(c.Key)
	if err != nil {
		logrus.Errorf("failed to get key when convert : %v", err)
	}
	return map[string]types.AttributeValue{"key": key}
}
