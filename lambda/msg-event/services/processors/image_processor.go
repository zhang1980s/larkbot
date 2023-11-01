package processors

import (
	"encoding/json"
	"msg-event/config"
	"msg-event/dao"
	"msg-event/model"
	"msg-event/model/event"
	"msg-event/services/api"
	"msg-event/utils"

	"github.com/sirupsen/logrus"
)

type imageProcessor struct {
}

func (r imageProcessor) ShouldProcess(e *event.Msg) bool {
	return true
}

func GetImageProcessor() api.Processor {
	return &imageProcessor{}
}

func (r imageProcessor) Process(e *event.Msg) error {
	c, err := dao.GetCaseByEvent(e)
	if err != nil {
		return err
	}
	content := &model.Content{}
	if err = json.Unmarshal([]byte(e.Event.Message.Content), content); err != nil {
		return err
	}
	data, err := dao.DownloadImage(e.Event.Message.MsgID, content.ImageKey)
	if err != nil {
		return err
	}
	format := utils.GuessImageFormat(data)
	err = dao.AddAttachmentToCase(c, content.ImageKey+format, data)
	if err != nil {
		logrus.Errorf("failed to att attachment %v", err)
		return err
	}

	dao.SendMsg(c.ChannelID, c.UserID, config.Conf.Ack)
	return nil
}
