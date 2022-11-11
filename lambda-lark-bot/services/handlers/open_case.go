package handlers

import (
	"lambda-lark-bot/config"
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"

	"github.com/sirupsen/logrus"
)

type openCaseServ struct {
}

func GetOpenCaseServ() api.Server {
	return &openCaseServ{}
}

func (s *openCaseServ) Handle(e *event.Msg, title string) (c *dao.Case, err error) {
	fromChannelID := e.Event.Message.ChatID
	customerID := e.Event.Sender.SenderIDs.UserID
	config.Conf.CaseCardTemplate.ChatId = fromChannelID
	config.Conf.CaseCardTemplate.UserId = customerID

	c = &dao.Case{
		Title: title,
	}
	rsp, err := dao.SendCardMsg(config.Conf.CaseCardTemplate, c)
	if err != nil {
		logrus.Errorf("Failed to send card msg, %v", err)
		return nil, err
	}
	return dao.OpenCase(fromChannelID, customerID, title, rsp.Data.MsgID, config.Conf.CaseCardTemplate)

}
