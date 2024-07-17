package handlers

import (
	"context"
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/services/api"
)

type qServ struct {
}

func GetQService() api.Server {
	return &qServ{}
}

func (s *qServ) Handle(e *event.Msg, content string) (c *dao.Case, err error) {

	dao.SendContentToSQS(context.Background(), content, e.Event.Message.MsgID)
	// fromChannelID := e.Event.Message.ChatID
	// customerID := e.Event.Sender.SenderIDs.UserID
	// msgID := e.Event.Message.MsgID

	// dao.SendMsg(fromChannelID, customerID, content+" "+msgID)
	return nil, nil
}

func (s *qServ) ShouldHandle(e *event.Msg) bool {
	return true
}
