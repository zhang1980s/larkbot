package handlers

import (
	"case-refresh/dao"
	"case-refresh/model/event"
	"case-refresh/services/api"

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
	caze.CardRespMsgID = *rsp.Data.MessageId
	return dao.UpsertCase(caze)
}
