package api

import (
	"lambda-lark-bot/dao"
	"lambda-lark-bot/model/event"
)

type Server interface {
	Handle(e *event.Msg, str string) (c *dao.Case, err error)
}
