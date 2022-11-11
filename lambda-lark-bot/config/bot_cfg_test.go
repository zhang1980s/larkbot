package config

import (
	"lambda-lark-bot/model"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

func Test_JsonGen(t *testing.T) {
	c := Config{
		Key: "test",
	}
	c.Usage = Usage
	c.ServiceMap = ServiceMap
	c.SevMap = SevMap
	c.CaseCardTemplate = &model.FeiShuMsg{
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
					Content: "",
					Href:    model.Href{},
				},
			},
		},
	}
	t.Run("test", func(t *testing.T) {
		item, _ := dynamodbattribute.MarshalMap(c)
		logrus.Info(item)
	})
}
