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

type CreateChatTabsReq struct {
	ChatTabs *[]ChatTabs `json:"chat_tabs"`
}

type ChatTabs struct {
	TabName    string      `json:"tab_name,omitempty"`
	TabType    string      `json:"tab_type"`
	TabContent *TabContent `json:"tab_content,omitempty"`
	TabConfig  *TabConfig  `json:"tab_config,omitempty"`
}

type TabContent struct {
	URL           string `json:"url,omitempty"`
	Doc           string `json:"doc,omitempty"`
	MeetingMinute string `json:"meeting_minute,omitempty"`
}

type TabConfig struct {
	IconKey   string `json:"icon_key,omitempty"`
	IsBuiltIn bool   `json:"is_built_in,omitempty"`
}
