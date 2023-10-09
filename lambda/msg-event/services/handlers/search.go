package handlers

import (
	"fmt"
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/services/api"

	"github.com/sirupsen/logrus"
)

type searcher struct {
}

func GetSearcher() api.Server {
	return &searcher{}
}

func (s *searcher) Handle(e *event.Msg, time string) (c *dao.Case, err error) {
	fromChannelID := e.Event.Message.ChatID

	// // search by time
	cs, err := dao.GetCasesByTime(time)
	if err != nil {
		logrus.Errorf("Failed to search case, %v", err)
		return nil, err
	}

	title := ""
	for _, v := range cs {
		s := fmt.Sprintf("[%s](%s) %s \\n", v.CaseID, v.CaseURL, v.Title)
		title += s
	}
	_, err = dao.SendMsgToChannel(fromChannelID, title)
	if err != nil {
		logrus.Errorf("Failed to send msg for search case, %v", err)
		return nil, err
	}
	return nil, nil
}
