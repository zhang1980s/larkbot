package processors

import (
	"errors"
	"lambda-lark-bot/model/event"
	"lambda-lark-bot/services/api"

	"github.com/sirupsen/logrus"
)

type cardProcessor struct {
}

func GetCardProcessor() api.Processor {
	return &cardProcessor{}
}

func (r cardProcessor) Process(e *event.Msg) (err error) {
	if e.Action != nil {
		if v, ok := serverManager[e.Action.Value.Key]; ok {
			logrus.Infof("commond %s. value %s", e.Action.Value, e.Action.Option)
			_, err = v.Handle(e, e.Action.Option)
			if err != nil {
				logrus.Errorf("faile to handle card msg %v", err)
				return err
			}
			return nil
		} else {
			logrus.Errorf("card select failed %v", e.Action)
			return errors.New("failed to match action handler")
		}
	}
	return nil
}
