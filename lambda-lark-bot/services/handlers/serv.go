package handlers

import (
	"lambda-lark-bot/config"
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type serv struct {
}

func GetServ() api.Server {
	return &serv{}
}

func (s *serv) Handle(e *event.Msg, str string) (c *dao.Case, err error) {

	c, err = dao.GetCaseByEvent(e)
	if err != nil {
		logrus.Errorf("failed to get case %s", err)
		return nil, err
	}
	service := strings.Trim(str, " ")
	_, ok := config.SevMap[service]
	if !ok {
		c.SevCode = "normal"
	} else {
		c.SevCode = service
	}

	c.UpdateTime = time.Now().String()
	for i, element := range c.CardMsg.Card.Elements {
		if element.Extra.Value.Key == e.Action.Value.Key {
			c.CardMsg.Card.Elements[i].Extra.InitialOption = e.Action.Option
		}
	}
	return dao.UpsertCase(c)
}
