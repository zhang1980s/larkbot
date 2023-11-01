package processors

import (
	"encoding/json"
	"msg-event/config"
	"msg-event/dao"
	"msg-event/model"
	"msg-event/model/event"
	"msg-event/services/api"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type textProcessor struct {
}

func GetTextProcessor() api.Processor {
	return &textProcessor{}
}

func (r textProcessor) ShouldProcess(e *event.Msg) bool {
	//perimission judgetment
	userId := e.Event.Sender.SenderIDs.UserID
	_, ok := config.Conf.UserWhiteListMap[userId]

	if os.Getenv("ENABLE_USER_WHITELIST") == "true" && !ok {
		fromChannelID := e.Event.Message.ChatID
		dao.SendMsgToChannel(fromChannelID, config.Conf.NoPermissionMSG)
		return false
	}
	return true
}

func (r textProcessor) Process(e *event.Msg) (err error) {
	if e.Event.Message.MsgType == "text" {
		if ok := r.ShouldProcess(e); !ok {
			return nil
		}

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
	return err
}
