package processors

import (
	"case-refresh/dao"
	"case-refresh/model"
	"case-refresh/model/event"
	"case-refresh/services/api"
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

type textProcessor struct {
}

func GetTextProcessor() api.Processor {
	return &textProcessor{}
}

func (r textProcessor) Process(e *event.Msg) (err error) {
	if e.Event.Message.MsgType == "text" {
		c := &model.Content{}
		if err = json.Unmarshal([]byte(e.Event.Message.Content), c); err != nil {
			return err
		}
		tokens := strings.SplitN(strings.Trim(c.Text, " "), " ", 2)
		cmd := tokens[0]
		content := ""
		if len(tokens) == 2 {
			content = tokens[1]
		}
		logrus.Infof("cmd %s, rest %s", cmd, content)

		if v, ok := serverManager[cmd]; ok {
			logrus.Infof("commond %s. content %s", cmd, content)
			_, err = v.Handle(e, content)
		} else {
			logrus.Infof("default as case comment %s", c.Text)
			v = serverManager[defaultKey]
			_, err = v.Handle(e, c.Text)
		}
		if err != nil {
			logrus.Errorf("process case failed %v", err)
			dao.SendErrCardMsg(e.Event.Message.ChatID, e.Event.Sender.SenderIDs.UserID, err)
		}

	}
	return nil
}
