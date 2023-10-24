package handlers

import (
	"fmt"
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/services/api"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type WhitlistServ struct {
}

func GetWhistlist() api.Server {
	return &WhitlistServ{}
}

func (s *WhitlistServ) Handle(e *event.Msg, whitelist string) (c *dao.Case, err error) {
	fromChannelID := e.Event.Message.ChatID
	// // get operation: add/del/list
	whitelistItem := strings.Split(whitelist, ",")

	var emailList []string
	var phoneList []string

	for _, item := range whitelistItem {

		rtn := isEmail(item)
		if rtn == true {
			emailList = append(emailList, item)
		} else {
			phoneList = append(phoneList, item)
		}
	}

	// //get userid by whitelist Email
	var badUserList []string
	var msg string
	validUer, badUserList, err := dao.GetUserIdbyEmailOrPhone(emailList, phoneList)
	if err != nil {
		msg = "无法获取用户id"
		_, err = dao.SendMsgToChannel(fromChannelID, msg)
		if err != nil {
			logrus.Errorf("Failed to send msg for whitelist, %v", err)
		}
		return nil, nil
	}

	if len(badUserList) > 0 {
		msg = fmt.Sprintf("无法获取 %v 对应的用户id，请核对是否是正确的电话号码或者邮箱", badUserList)
		_, err = dao.SendMsgToChannel(fromChannelID, msg)
		if err != nil {
			logrus.Errorf("Failed to send msg for whitelist, %v", err)
		}
		return nil, nil
	}

	// //write userID to ddb
	err = dao.AddWhitelist(validUer)
	if err != nil {
		msg = "添加白名单失败，请重试"
	} else {
		msg = "添加白名单成功"
	}

	_, err = dao.SendMsgToChannel(fromChannelID, msg)
	return nil, nil
}

func (s *WhitlistServ) ShouldHandle(e *event.Msg) bool {
	return true
}

func isEmail(s string) bool {
	// 这是一个简单的邮箱正则，根据需要可以进一步完善
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(s)
}
