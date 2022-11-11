package config

import (
	"lambda-lark-bot/model"
)

const (
	CaseNotExisted      = "您还没有开工单，请按照说明开工单"
	BotConfigNotExisted = "机器人配置丢失，并且创建默认失败"
)

//Usage text
var Usage = `开AWS支持案例方法：@机器人 开工单 问题
问题 [工单的题目]
内容 [工单内容] [包括：问题发生的时间及时区/涉及的资源ID及region/发生问题的现象/该问题对业务造成的影响/联系人及联系方式等信息]

账号 [问题涉及资源属于的AWS账户]
响应速度 [low - 24小时 normal - 12小时 high - 4小时 urgent - 1小时 critical - 15分钟) ]
服务 [问题涉及服务]

案例更新：[在机器人创建的新工单群里发言提交工单更新]`

var ErrCardTemplate = &model.FeiShuMsg{
	ChatId:      "",
	MsgType:     "",
	UpdateMulti: false,
	Card: model.Card{
		Config: model.Config{
			WideScreenMode: true,
		},
		Elements: []model.Elements{
			{
				Tag: "markdown",
			},
		},
	},
}

var CardTemplate = &model.FeiShuMsg{
	UserId:      "test",
	Email:       "test",
	MsgType:     "test",
	Content:     nil,
	OpenId:      "test",
	RootId:      "test",
	ChatId:      "test",
	UpdateMulti: false,
	Card: model.Card{
		Config: model.Config{
			WideScreenMode: true,
		},
		Elements: []model.Elements{
			{
				Tag: "div",
				Text: model.Text{
					Tag:     "lark_md",
					Content: "级别",
				},
				Extra: model.Extra{
					Tag: "select_static",
					Placeholder: model.Placeholder{
						Tag:     "plain_text",
						Content: "默认提示",
					},
					Value: model.Value{
						Key: "key1",
					},
					Options: []model.Options{
						{
							Text: model.Text{
								Tag:     "plain_text",
								Content: "选项一",
							},
							Value: "1",
						},
					},
				},
			},
		},
	},
}

var SevMap = map[string]string{
	"low":      "low",
	"normal":   "normal",
	"high":     "high",
	"urgent":   "urgent",
	"critical": "critical",
}

var ServiceMap = map[string][]string{
	"0":  {"general-info", "using-aws"},
	"1":  {"amazon-elastic-compute-cloud-linux", "other"},
	"2":  {"amazon-simple-storage-service", "general-guidance"},
	"3":  {"amazon-virtual-private-cloud", "general-guidance"},
	"4":  {"elastic-load-balancing", "general-guidance"},
	"5":  {"aws-identity-and-access-management", "general-guidance"},
	"6":  {"amazon-cloudwatch", "general-guidance"},
	"7":  {"aws-direct-connect", "general-guidance"},
	"8":  {"distributed-denial-of-service", "inbound-to-aws"},
	"9":  {"account-management", "billing"},
	"10": {"amazon-cloudfront", "general-guidance"},
	"11": {"amazon-relational-database-service-postgresql", "general-guidance"},
	"12": {"amazon-relational-database-service-mysql", "general-guidance"},
	"13": {"amazon-elastic-block-store", "general-guidance"},
	"14": {"aws-lambda", "general-guidance"},
}
