package processors

import (
	"encoding/json"
	"msg-event/config"
	"msg-event/dao"
	"msg-event/model"
	"msg-event/model/event"
	"msg-event/services/api"

	"github.com/sirupsen/logrus"
)

type attaProcessor struct {
}

func (r attaProcessor) ShouldProcess(e *event.Msg) bool {
	return true
}

func GetAttaProcessor() api.Processor {
	return &attaProcessor{}
}

func (r attaProcessor) Process(e *event.Msg) error {
	c, err := dao.GetCaseByEvent(e)
	if err != nil {
		return err
	}

	content := &model.Content{}
	if err = json.Unmarshal([]byte(e.Event.Message.Content), content); err != nil {
		return err
	}
	data, err := dao.DownloadFile(e.Event.Message.MsgID, content.FileKey)
	if err != nil {
		return err
	}
	err = dao.AddAttachmentToCase(c, content.FileName, data)
	if err != nil {
		logrus.Errorf("failed to att attachment %v", err)
		return err
	}

	dao.SendMsg(c.ChannelID, c.UserID, config.Conf.Ack)

	return nil
}
