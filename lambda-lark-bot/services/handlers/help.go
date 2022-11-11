package handlers

import (
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"

	"github.com/sirupsen/logrus"
)

type helper struct {
}

func Gethelper() api.Server {
	return &helper{}
}

func (h *helper) Handle(e *event.Msg, title string) (c *dao.Case, err error) {
	caze, err := dao.GetCaseByEvent(e)
	if err != nil {
		return nil, err
	}
	rsp, err := dao.SendCardMsg(caze.CardMsg, caze)
	if err != nil {
		logrus.Errorf("send card msg failed, %v", err)
	}
	caze.CardRespMsgID = rsp.Data.MsgID
	return dao.UpsertCase(caze)
}
