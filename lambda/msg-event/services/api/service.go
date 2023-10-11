package api

import (
	"msg-event/dao"
	"msg-event/model/event"
)

type Server interface {
	Handle(e *event.Msg, str string) (c *dao.Case, err error)
	ShouldHandle() bool
}
