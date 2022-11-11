package config

import "lambda-lark-bot/model"

var Conf *Config

type Config struct {
	Key              string              `json:"key,omitempty"`
	SevMap           map[string]string   `json:"sev_map,omitempty"`
	ServiceMap       map[string][]string `json:"service_map,omitempty"`
	Usage            string              `json:"usage,omitempty"`
	Accounts         map[string]*Account `json:"accounts,omitempty"`
	AppID            string              `json:"app_id,omitempty"`
	AppSecret        string              `json:"app_secret,omitempty"`
	ErrCardTemplate  *model.FeiShuMsg    `json:"err_card_template,omitempty"`
	CaseCardTemplate *model.FeiShuMsg    `json:"case_card_template,omitempty"`
	Ack              string              `json:"ack,omitempty"`
}

type Account struct {
	AccessKeyID     string `json:"access_key_id,omitempty"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
}
