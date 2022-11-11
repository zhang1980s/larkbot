package response

import "lambda-lark-bot/model"

type MsgResponse struct {
	Challenge string           `json:"challenge,omitempty"`
	Elements  []model.Elements `json:"elements,omitempty"`
}
