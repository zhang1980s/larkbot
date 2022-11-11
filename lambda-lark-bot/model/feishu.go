package model

type FeiShuResponse struct {
	Code int    `json:"code,omitempty"`
	Data *Data  `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
type Data struct {
	ChatID string `json:"chat_id,omitempty"`
}

type TokenReq struct {
	AppID     string `json:"app_id,omitempty"`
	AppSecret string `json:"app_secret,omitempty"`
}

type TokenResp struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	TAToken string `json:"tenant_access_token,omitempty"`
}

type Content struct {
	Text     string `json:"text,omitempty"`
	ImageKey string `json:"image_key,omitempty"`
	FileKey  string `json:"file_key,omitempty"`
	FileName string `json:"file_name,omitempty"`
}

type FeiShuMsg struct {
	UserId      string   `json:"user_id,omitempty"`
	Email       string   `json:"email,omitempty"`
	MsgType     string   `json:"msg_type,omitempty"`
	Content     *Content `json:"content,omitempty"`
	OpenId      string   `json:"open_id,omitempty"`
	RootId      string   `json:"root_id,omitempty"`
	ChatId      string   `json:"chat_id,omitempty"`
	UpdateMulti bool     `json:"update_multi"`
	Card        Card     `json:"card"`
}

type I18nNames struct {
	ZhCn string `json:"zh_cn,omitempty"`
	EnUs string `json:"en_us,omitempty"`
}

type CreateChannelReq struct {
	UserIds     []string   `json:"user_ids,omitempty"`
	OpenIds     []string   `json:"open_ids,omitempty"`
	I18nNames   *I18nNames `json:"i18n_names,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
}

type CardRsp struct {
	Data CRData `json:"data,omitempty"`
}

type CRData struct {
	MsgID string `json:"message_id,omitempty"`
}
