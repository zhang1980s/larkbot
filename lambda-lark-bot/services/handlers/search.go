package handlers

import (
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"
)

type searcher struct {
}

func GetSearcher() api.Server {
	return &searcher{}
}

func (s *searcher) Handle(e *event.Msg, time string) (c *dao.Case, err error) {
	// fromChannelID := e.Event.Message.ChatID
	// customerID := e.Event.Sender.SenderIDs.UserID

	// // search by time
	// cs, err = dao.GetCasesByTime(time)
	// if err != nil {
	// 	logrus.Errorf("Failed to search case, %v", err)
	// 	return nil, err
	// }
	// title := ""
	// for _, v := range cs {
	// 	title += v.CaseID + v.Title
	// }
	// err = dao.SendMsg(fromChannelID, customerID, title)
	// if err != nil{
	// 	logrus.Errorf("Failed to send msg for search case, %v", err)
	// 	return nil, err
	// }
	return nil, nil
}
