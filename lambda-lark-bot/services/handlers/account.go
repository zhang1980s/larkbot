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

type accountServ struct {
}

func GetAccountServ() api.Server {
	return &accountServ{}
}

func (s *accountServ) Handle(e *event.Msg, str string) (c *dao.Case, err error) {
	c, err = dao.GetCaseByEvent(e)
	if err != nil {
		logrus.Errorf("failed to get case %s", err)
		return nil, err
	}
	account := strings.Trim(str, " ")
	_, ok := config.Conf.Accounts[account]
	if !ok {
		c.AccountKey = "0"
	} else {
		c.AccountKey = account
	}
	c.UpdateTime = time.Now().String()
	for i, element := range c.CardMsg.Card.Elements {
		if element.Extra.Value.Key == e.Action.Value.Key {
			c.CardMsg.Card.Elements[i].Extra.InitialOption = e.Action.Option
		}
	}
	return dao.UpsertCase(c)
}
