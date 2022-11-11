package processors

import (
	"encoding/json"
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"
	"lambda-lark-bot/utils"

	"github.com/sirupsen/logrus"
)

type imageProcessor struct {
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
	return nil
}
