package services

import (
	"context"
	"errors"
	"msg-event/config"
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/model/response"
	"msg-event/services/api"
	"msg-event/services/processors"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var processorManager map[string]api.Processor

func InitProcessors() {
	processorManager = map[string]api.Processor{
		"fresh_comment": processors.GetRefreshCommentProcessor(),
		"card":          processors.GetCardProcessor(),
		"text":          processors.GetTextProcessor(),
		"image":         processors.GetImageProcessor(),
		"file":          processors.GetAttaProcessor(),
	}
}

func Serve(_ context.Context, e *event.Msg) (event *response.MsgResponse, err error) {
	logrus.Infof("====================================================")
	err = dao.SetupConfig()
	processors.InitServices()
	InitProcessors()
	resp := &response.MsgResponse{
		Challenge: e.Challenge,
	}
	if err != nil {
		logrus.Errorf("setup config failed %s", err)
		return resp, err
	}

	logrus.Infof("USER ID: %v", e.Event.Sender.SenderIDs.UserID)
	_, ok := config.Conf.UserWhiteListMap[e.Event.Sender.SenderIDs.UserID]
	if os.Getenv("ENABLE_WL") == "true" && !ok {
		fromChannelID := e.Event.Message.ChatID
		dao.SendMsgToChannel(fromChannelID, config.Conf.NoPermissionMSG)
		return resp, nil
	}

	if e.Action != nil && e.Event.Message.MsgType == "" {
		e.Event.Message.MsgType = "card"
	}
	if e.Event.Message.MsgType != "" {
		if !Processable(e) {
			logrus.Infof("Duplicate message with same eventID")
			return resp, nil
		}
		if v, ok := processorManager[e.Event.Message.MsgType]; ok {
			logrus.Infof("event type %s. ", e.Event.Message.MsgType)

			err = v.Process(e)
			if err != nil {
				logrus.Errorf("failed to process %v", err)
				return resp, err
			}
		} else {
			logrus.Errorf("unknown type %s", e.Event.Message.MsgType)
			return resp, err
		}
	}

	caze, err := dao.GetCaseByEvent(e)
	if err != nil {
		logrus.Errorf("failed to get case, %v", err)
		return resp, err
	}
	if caze == nil {
		logrus.Infof("Return challenge for url_verification")
		return resp, nil
	}

	if caze != nil && caze.Status == dao.STATUS_NEW {

		// if user in list, create case and channel, else send no permission

		err = createChatOrNewCase(caze)
		if err != nil {
			logrus.Errorf("process chat or create case failed case %+v, \n %v", caze, err)
		}

	}
	resp.Elements = caze.CardMsg.Card.Elements
	return resp, nil
}

func createChatOrNewCase(caze *dao.Case) error {
	caze.Print()
	_, sevOk := config.SevMap[caze.SevCode]
	_, serviceOk := config.ServiceMap[caze.ServiceCode]
	_, accountOK := config.Conf.Accounts[caze.AccountKey]

	if strings.Trim(caze.Title, " ") != "" &&
		strings.Trim(caze.Content, " ") != "" &&
		strings.Trim(caze.SevCode, " ") != "" &&
		sevOk &&
		strings.Trim(caze.ServiceCode, " ") != "" &&
		serviceOk &&
		accountOK &&
		caze.Status == dao.STATUS_NEW {

		caze.Status = dao.STATUS_OPEN
		caze, err := dao.CreateCaseAndChannel(caze)
		if err != nil {
			logrus.Errorf("failed to create case info %s", err)
			return err
		}
		//clean up fromchannel
		caze.ChannelID = caze.FromChannelID
		caze.UserID = ""
		caze.Title = ""
		caze.Content = ""
		caze.Type = dao.TYPE_OPEN_CASE
		caze.SevCode = ""
		caze.ServiceCode = ""
		caze.CardMsg.ChatId = caze.ChannelID
		caze.CardMsg.UserId = caze.UserID
		_, err = dao.UpsertCase(caze)
		if err != nil {
			logrus.Errorf("failed to cleanup root case info %s", err)
			return err
		}
		return nil
	}
	s := dao.FormatMsg(caze)
	return errors.New(s)
}
