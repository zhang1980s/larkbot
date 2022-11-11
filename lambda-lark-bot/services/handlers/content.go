package handlers

import (
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"
	"time"

	"github.com/sirupsen/logrus"
)

type contentServ struct {
}

func GetContentServ() api.Server {
	return &contentServ{}
}

func (s *contentServ) Handle(e *event.Msg, content string) (c *dao.Case, err error) {
	c, err = dao.GetCaseByEvent(e)
	if err != nil {
		return nil, err
	}

	c.Content = content
	c.UpdateTime = time.Now().String()

	rsp, err := dao.SendCardMsg(c.CardMsg, c)
	if err != nil {
		logrus.Errorf("send card msg failed, %v", err)
	}
	c.CardRespMsgID = rsp.Data.MsgID
	return dao.UpsertCase(c)
}
