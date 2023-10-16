package handlers

import (
	"fmt"
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/services/api"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
)

type searcher struct {
}

func GetSearcher() api.Server {
	return &searcher{}
}

func (s *searcher) Handle(e *event.Msg, time string) (c *dao.Case, err error) {
	fromChannelID := e.Event.Message.ChatID

	// // search by time
	cs, err := dao.GetCasesByTime(time)
	if err != nil {
		logrus.Errorf("Failed to search case, %v", err)
		return nil, err
	}

	title := "案例号\\t\\t\\t 账号\\t\\t\\t 创建时间\\t\\t\\t\\t 案例状态\\t\\t 题目 \\n"
	for _, v := range cs {
		// if v.Status != "NEW" && v.Status != "OPEN" {
		if v.Status != "NEW" && v.Title != "" {
			s := fmt.Sprintf("[%s](%s)\\t %s\\t %s\\t %s\\t\\t %s\\n", v.DisplayCaseID, v.CaseURL, v.CaseAccountID, formatTimestype(v.CreateTime), v.Status, v.Title)
			title += s
		}

	}
	_, err = dao.SendMsgToChannel(fromChannelID, title)
	if err != nil {
		logrus.Errorf("Failed to send msg for search case, %v", err)
		return nil, err
	}
	return nil, nil
}

func formatTimestype(input string) string {
	inputFormat := "2006-01-02 15:04:05.999999999 -0700 MST"
	outputFormat := "January 2, 2006 MST"

	re := regexp.MustCompile(` m=[+-]?\d+\.\d+`)

	input = re.ReplaceAllString(input, "")

	t, err := time.Parse(inputFormat, input)
	if err != nil {
		logrus.Infof("Error parsing input: %v", err)
		return ""
	}

	output := t.Format(outputFormat)

	return output
}

func (s *searcher) ShouldHandle(e *event.Msg) bool {
	return true
}
