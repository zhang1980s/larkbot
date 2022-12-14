package processors

import (
	"encoding/json"
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"

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
	return nil
}
