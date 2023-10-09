package api

import (
	"case-refresh/dao"
	"case-refresh/model/event"
)

type Server interface {
	Handle(e *event.Msg, str string) (c *dao.Case, err error)
}
