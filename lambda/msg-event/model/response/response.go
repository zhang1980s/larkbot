package response

import "msg-event/model"

type MsgResponse struct {
	Challenge string           `json:"challenge,omitempty"`
	Elements  []model.Elements `json:"elements,omitempty"`
}
