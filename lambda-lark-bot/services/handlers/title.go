package handlers

import (
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"
	"time"

	"github.com/sirupsen/logrus"
)

type titleServ struct {
}

func GetTitleServ() api.Server {
	return &titleServ{}
}

func (s *titleServ) Handle(e *event.Msg, title string) (c *dao.Case, err error) {
	c, err = dao.GetCaseByEvent(e)
	if err != nil {
		return nil, err
	}

	c.Title = title
	c.UpdateTime = time.Now().String()

	rsp, err := dao.SendCardMsg(c.CardMsg, c)
	if err != nil {
		logrus.Errorf("send card msg failed, %v", err)
	}
	c.CardRespMsgID = rsp.Data.MsgID
	return dao.UpsertCase(c)
}
