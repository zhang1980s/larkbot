package handlers

import (
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"

	"github.com/sirupsen/logrus"
)

type searcher struct {
}

func GetSearcher() api.Server {
	return &searcher{}
}

func (s *searcher) Handle(e *event.Msg, title string) (c *dao.Case, err error) {
	/// print test point
	logrus.Infof("%v", c)
	logrus.Infof("Into seach loop %s", title)

	dao.SendMsg(c.ChannelID, c.UserID, title)
	return nil, nil
}
