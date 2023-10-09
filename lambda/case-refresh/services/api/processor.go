package api

import "case-refresh/model/event"

type Processor interface {
	Process(e *event.Msg) error
}
