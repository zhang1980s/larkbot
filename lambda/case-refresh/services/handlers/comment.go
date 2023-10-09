package handlers

import (
	"case-refresh/config"
	"case-refresh/dao"
	"case-refresh/model/event"
	"case-refresh/services/api"
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
)

type commentsServ struct {
}

func GetCommentsServServ() api.Server {
	return &commentsServ{}
}

func (s *commentsServ) Handle(e *event.Msg, str string) (c *dao.Case, err error) {
	c, err = dao.GetCaseByEvent(e)
	if err != nil {
		logrus.Errorf("get case failed %+v", err)
		return nil, errors.New(config.CaseNotExisted)
	}
	cazeID := strings.Trim(c.CaseID, " ")
	if cazeID == "" || c.Type == dao.TYPE_OPEN_CASE {
		return nil, errors.New(dao.FormatMsg(c))
	}

	// add comment to aws case system
	c, err = dao.AddComment(c, str)
	if err != nil {
		logrus.Errorf("add comment failed %+v", err)
		return nil, err
	}

	dao.SendMsg(c.ChannelID, c.UserID, config.Conf.Ack)
	return c, nil
}
