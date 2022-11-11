package api

import "lambda-lark-bot/model/event"

type Processor interface {
	Process(e *event.Msg) error
}
