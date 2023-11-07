package handlers

import (
	"errors"
	"fmt"
	"msg-event/config"
	"msg-event/dao"
	"msg-event/model/event"
	"msg-event/services/api"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type WhitlistServ struct {
}

type WhitelistDelServ struct {
}

type WhitelistCatServ struct {
}

type AdminWhitelistServ struct {
}

func GetWhistlist() api.Server {
	return &WhitlistServ{}
}

func GetWhitelistDel() api.Server {
	return &WhitelistDelServ{}
}

func GetWhitelistCat() api.Server {
	return &WhitelistCatServ{}
}

func GetAdminWhitelist() api.Server {
	return &AdminWhitelistServ{}
}

func (s *WhitlistServ) Handle(e *event.Msg, whitelist string) (c *dao.Case, err error) {
	if _, ok := config.Conf.RoleMap[e.Event.Sender.SenderIDs.UserID]; !ok {
		err = errors.New("你没有权限添加白名单")
		return nil, err
	}

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
		logrus.Errorf("Failed to send msg for whitelist, %v", err)
		return nil, err
	}

	if len(badUserList) > 0 {
		msg = fmt.Sprintf("无法获取 %v 对应的用户id，请核对是否是正确的电话号码或者邮箱", badUserList)
		err = errors.New(msg)
		return nil, err
	}

	// //write userID to ddb
	err = dao.AddWhitelist(validUer)
	if err != nil {
		msg = "添加白名单失败，请重试"
		err = errors.New(msg)
		return nil, err
	}
	msg = "添加白名单成功"
	fromChannelID := e.Event.Message.ChatID
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

func (s *WhitelistDelServ) ShouldHandle(e *event.Msg) bool {
	return true
}

func (s *WhitelistCatServ) ShouldHandle(e *event.Msg) bool {
	return true
}

func (s *AdminWhitelistServ) ShouldHandle(e *event.Msg) bool {
	return true
}

func (s *WhitelistDelServ) Handle(e *event.Msg, whitelist string) (c *dao.Case, err error) {
	if _, ok := config.Conf.RoleMap[e.Event.Sender.SenderIDs.UserID]; !ok {
		err = errors.New("你没有权限删除白名单")
		return nil, err
	}

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

	var badUserList []string
	var msg string
	validUer, badUserList, err := dao.GetUserIdbyEmailOrPhone(emailList, phoneList)
	if err != nil {
		logrus.Errorf("Failed to send msg for whitelist, %v", err)
		return nil, err
	}

	if len(badUserList) > 0 {
		msg = fmt.Sprintf("无法获取 %v 对应的用户id，请核对是否是正确的电话号码或者邮箱", badUserList)
		err = errors.New(msg)
		return nil, err
	}

	err = dao.DelWhiteList(validUer)
	if err != nil {
		msg = "删除白名单失败，请重试"
		err = errors.New(msg)
		return nil, err
	}
	msg = "删除白名单成功"
	fromChannelID := e.Event.Message.ChatID
	_, err = dao.SendMsgToChannel(fromChannelID, msg)
	return nil, nil
}

func (s *WhitelistCatServ) Handle(e *event.Msg, whitelist string) (c *dao.Case, err error) {
	if _, ok := config.Conf.RoleMap[e.Event.Sender.SenderIDs.UserID]; !ok {
		err = errors.New("你没有权限查看白名单")
		return nil, err
	}
	rtnWhitelist := dao.GetWhiteList()
	msg := "白名单: "
	for user, role := range rtnWhitelist {
		msg += fmt.Sprintf("%s:%s ;", user, role)
	}
	logrus.Info(msg)
	fromChannelID := e.Event.Message.ChatID
	_, err = dao.SendMsgToChannel(fromChannelID, msg)

	return nil, nil
}

func (s *AdminWhitelistServ) Handle(e *event.Msg, whitelist string) (c *dao.Case, err error) {
	if _, ok := config.Conf.RoleMap[e.Event.Sender.SenderIDs.UserID]; !ok {
		err = errors.New("你没有权限设置管理员")
		return nil, err
	}

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

	var badUserList []string
	var msg string
	validUer, badUserList, err := dao.GetUserIdbyEmailOrPhone(emailList, phoneList)
	if err != nil {
		logrus.Errorf("Failed to send msg for whitelist, %v", err)
		return nil, err
	}

	if len(badUserList) > 0 {
		msg = fmt.Sprintf("无法获取 %v 对应的用户id，请核对是否是正确的电话号码或者邮箱", badUserList)
		err = errors.New(msg)
		return nil, err
	}
	err = dao.SetAdmin(validUer)
	if err != nil {
		msg = "添加admin失败，请重试"
		err = errors.New(msg)
		return nil, err
	}
	msg = "添加admin成功"
	fromChannelID := e.Event.Message.ChatID
	_, err = dao.SendMsgToChannel(fromChannelID, msg)

	return nil, nil
}
