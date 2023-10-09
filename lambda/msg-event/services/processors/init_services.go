package processors

import (
	"msg-event/services/api"
	"msg-event/services/handlers"
)

var serverManager map[string]api.Server
var defaultKey = "default"

func InitServices() {
	serverManager = map[string]api.Server{
		"开工单":      handlers.GetOpenCaseServ(),
		"内容":       handlers.GetContentServ(),
		"账户":       handlers.GetAccountServ(),
		"问题":       handlers.GetTitleServ(),
		"响应速度":     handlers.GetServ(),
		"服务":       handlers.GetServiceServ(),
		"帮助":       handlers.Gethelper(),
		"历史":       handlers.GetSearcher(),
		defaultKey: handlers.GetCommentsServServ(),
	}
}
