package api

import "msg-event/model/event"

type Processor interface {
	Process(e *event.Msg) error
}
