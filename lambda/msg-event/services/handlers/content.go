package handlers

import (
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/services/api"
	"time"

	"github.com/sirupsen/logrus"
)

const contentKey = "content"

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

	for i, element := range c.CardMsg.Card.Elements {
		if element.Extra.Value.Key == contentKey {
			c.CardMsg.Card.Elements[i].Content += content
			break
		}
	}

	rsp, err := dao.SendCardMsg(c.CardMsg, c)
	if err != nil {
		logrus.Errorf("send card msg failed, %v", err)
	}
	c.CardRespMsgID = *rsp.Data.MessageId
	return dao.UpsertCase(c)
}
