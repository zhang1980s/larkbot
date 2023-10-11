package handlers

import (
	"msg-event/config"
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/services/api"

	"github.com/sirupsen/logrus"
)

const openCaseTitleKey = "title"

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

	for i, element := range config.Conf.CaseCardTemplate.Card.Elements {
		if element.Extra.Value.Key == openCaseTitleKey {
			config.Conf.CaseCardTemplate.Card.Elements[i].Content += title
			logrus.Infof("match key %v. value %v", openCaseTitleKey, title)
			break
		} else {
			logrus.Infof("not match key %v. value %v", openCaseTitleKey, title)
		}
	}

	rsp, err := dao.SendCardMsg(config.Conf.CaseCardTemplate, c)
	if err != nil {
		logrus.Errorf("Failed to send card msg, %v", err)
		return nil, err
	}
	return dao.OpenCase(fromChannelID, customerID, title, *rsp.Data.MessageId, config.Conf.CaseCardTemplate)

}

func (s *openCaseServ) ShouldHandle() bool {
	return true
}
