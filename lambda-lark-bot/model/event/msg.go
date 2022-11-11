package event

type Header struct {
	EventID    string `json:"event_id,omitempty"`
	EventType  string `json:"event_type,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
	Token      string `json:"token,omitempty"`
	AppID      string `json:"app_id,omitempty"`
	TenantKey  string `json:"tenant_key,omitempty"`
}

type Event struct {
	Sender  Sender  `json:"sender,omitempty"`
	Message Message `json:"message,omitempty"`
}

type Sender struct {
	SenderIDs  SenderIDS `json:"sender_id,omitempty"`
	SenderType string    `json:"sender_type,omitempty"`
	TenantKey  string    `json:"tenant_key,omitempty"`
}

type SenderIDS struct {
	UnionID string `json:"union_id,omitempty"`
	UserID  string `json:"user_id,omitempty"`
	OpenID  string `json:"sender_id,omitempty"`
}

type Message struct {
	MsgID      string `json:"message_id,omitempty"`
	RootID     string `json:"root_id,omitempty"`
	ParentID   string `json:"parent_id,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
	ChatID     string `json:"chat_id,omitempty"`
	ChatType   string `json:"chat_type,omitempty"`
	Content    string `json:"content,omitempty"`
	MsgType    string `json:"message_type,omitempty"`
}

type Msg struct {
	Schema    string `json:"schema,omitempty"`
	Event     Event  `json:"event,omitempty"`
	Challenge string `json:"challenge"`
	Header    Header `json:"header,omitempty"`
	//card
	OpenID    string  `json:"open_id"`
	UserID    string  `json:"user_id"`
	TenantKey string  `json:"tenant_key"`
	OpenMsgID string  `json:"open_message_id"`
	Token     string  `json:"token"`
	Action    *Action `json:"action"`
}

type Action struct {
	Value  *Value `json:"value"`
	Tag    string `json:"tag"`
	Option string `json:"option"`
}

type Value struct {
	Key string `json:"key"`
}
