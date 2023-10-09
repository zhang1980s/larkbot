package handlers

import (
	"msg-event/config"
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/services/api"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type serviceServ struct {
}

func GetServiceServ() api.Server {
	return &serviceServ{}
}

func (s *serviceServ) Handle(e *event.Msg, str string) (c *dao.Case, err error) {
	c, err = dao.GetCaseByEvent(e)
	if err != nil {
		logrus.Errorf("failed to get case %s", err)
		return nil, err
	}
	service := strings.Trim(str, " ")
	_, ok := config.ServiceMap[service]
	if !ok {
		c.ServiceCode = "1"
	} else {
		c.ServiceCode = service
	}
	c.UpdateTime = time.Now().String()
	for i, element := range c.CardMsg.Card.Elements {
		if element.Extra.Value.Key == e.Action.Value.Key {
			c.CardMsg.Card.Elements[i].Extra.InitialOption = e.Action.Option
		}
	}
	return dao.UpsertCase(c)
}
