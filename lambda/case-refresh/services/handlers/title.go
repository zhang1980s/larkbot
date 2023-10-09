package handlers

import (
	"case-refresh/dao"
	"case-refresh/model/event"
	"case-refresh/services/api"
	"time"

	"github.com/sirupsen/logrus"
)

const titleKey = "title"

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

	for i, element := range c.CardMsg.Card.Elements {
		if element.Extra.Value.Key == titleKey {
			c.CardMsg.Card.Elements[i].Content += title
			logrus.Infof("match key %v. value %v", titleKey, title)
			break
		} else {
			logrus.Infof("not match key %v. value %v", titleKey, title)
		}
	}
	rsp, err := dao.SendCardMsg(c.CardMsg, c)
	if err != nil {
		logrus.Errorf("send card msg failed, %v", err)
	}
	c.CardRespMsgID = *rsp.Data.MessageId
	return dao.UpsertCase(c)
}
