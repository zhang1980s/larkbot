package processors

import (
	"case-refresh/config"
	"case-refresh/dao"
	"case-refresh/model"
	"case-refresh/model/event"
	"case-refresh/services/api"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type attaProcessor struct {
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
