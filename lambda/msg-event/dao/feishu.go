package dao

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"msg-event/config"
	"msg-event/model"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	downloadUrl      = "https://open.feishu.cn/open-apis/im/v1/messages/%s/resources/%s?type=%s"
	tokenUrl         = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/"
	createChatTabUrl = "https://open.feishu.cn/open-apis/im/v1/chats/%s/chat_tabs"
)

func CreateChannel(userIDs []string, name string) (channelID string, err error) {
	client := getClient()

	req := larkim.NewCreateChatReqBuilder().
		UserIdType("user_id").
		Body(larkim.NewCreateChatReqBodyBuilder().
			UserIdList(userIDs).
			Name(name).
			Build()).
		Build()

	resp, err := client.Im.Chat.Create(context.Background(), req)
	if err != nil {
		logrus.Errorf("chat channel create failed, err %v", err)
		return "", err
	}
	if !resp.Success() {
		logrus.Errorf("chat channel create failed, response code %v", resp.Code)
		return "", errors.New(resp.CodeError.String())
	}

	logrus.Infof("response Body: %s", resp.RawBody)
	return *resp.Data.ChatId, nil
}

// SendMsg chatId group ID
// SendMsg userID userID
func SendMsg(chatId, userID, msg string) (resp *larkim.CreateMessageResp, err error) {
	TextMsg := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()
	return sendFeiShuMsg(getClient(), larkim.MsgTypeText, chatId, TextMsg)
}

func SendMsgToChannel(chatID, msg string) (resp *larkim.CreateMessageResp, err error) {
	return SendMsg(chatID, "", msg)
}

func getClient() *lark.Client {
	id, err := GetAppID()
	if err != nil {
		panic(err)
	}
	sec, err := GetAPPSecret()
	if err != nil {
		panic(err)
	}
	return lark.NewClient(id, sec)
}

func sendFeiShuMsg(client *lark.Client, t, chatId, msg string) (resp *larkim.CreateMessageResp, err error) {

	resp, err = client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(t).
			ReceiveId(chatId).
			Content(msg).
			Build()).
		Build())

	logrus.Infof("msg %v", msg)

	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		logrus.Infoln(resp.Code, resp.Msg, resp.RequestId())
		return resp, nil
	}

	logrus.Infoln(resp.Data.MessageId)
	logrus.Infoln(larkcore.Prettify(resp))
	logrus.Infoln(resp.RequestId())
	return resp, nil
}

func SendCardMsg(msgCard *model.FeiShuMsg, caze *Case) (*larkim.CreateMessageResp, error) {

	jsonStr, err := json.Marshal(msgCard.Card)
	if err != nil {
		return nil, err
	}
	resp, err := sendFeiShuMsg(getClient(), larkim.MsgTypeInteractive, msgCard.ChatId, string(jsonStr))
	if err != nil {
		logrus.Errorf("Failed to send card msg, %v", err)
		return nil, err
	}

	return resp, nil
}

func SendErrCardMsg(chatId, userID string, e error) error {
	config.Conf.ErrCardTemplate.Card.Elements[0].Content = e.Error()
	config.Conf.ErrCardTemplate.ChatId = chatId

	jsonStr, err := json.Marshal(config.Conf.ErrCardTemplate.Card)
	if err != nil {
		return err
	}
	resp, err := sendFeiShuMsg(getClient(), larkim.MsgTypeInteractive, chatId, string(jsonStr))
	if err != nil {
		logrus.Errorf("Failed to send card msg, %v", err)
		return err
	}
	logrus.Infof("Send err card rsp %v", string(resp.RawBody))
	return err
}

func getToken() (t *model.TokenResp, err error) {
	trq := &model.TokenReq{
		AppID:     config.Conf.AppID,
		AppSecret: config.Conf.AppSecret,
	}

	jsonStr, err := json.Marshal(trq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logrus.Infof("get token response Body: %s", string(body))
	t = &model.TokenResp{}
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func DownloadImage(msgID, imageKey string) ([]byte, error) {
	return download(fmt.Sprintf(downloadUrl, msgID, imageKey, "image"))
}

func DownloadFile(msgID, fileKey string) ([]byte, error) {
	return download(fmt.Sprintf(downloadUrl, msgID, fileKey, "file"))
}

func download(downloadUrl string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, downloadUrl, nil)
	if err != nil {
		return nil, err
	}
	t, err := getToken()
	if err != nil {
		logrus.Errorf("failed to get token %+v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+t.TAToken)

	logrus.Infof("downlaod req %v", req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}
	// Check server response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logrus.Infof("bad status: %s", string(body))
		return nil, fmt.Errorf("bad status: %v", resp)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("falied to read resp %v", err)
		return nil, err
	}
	return data, nil
}

func CreateChatTab(chatID string, url string) (err error) {
	ct := &model.CreateChatTabsReq{
		ChatTabs: &[]model.ChatTabs{
			{
				TabName: "CASELINK",
				TabType: "url",
				TabContent: &model.TabContent{
					URL: url,
				},
			},
		},
	}

	jsonStr, err := json.Marshal(ct)
	if err != nil {
		return err
	}

	/// Debug: print jsonStr
	///	logrus.Infof(string(jsonStr))

	createChatTabUrl := fmt.Sprintf(createChatTabUrl, chatID)
	/// logrus.Infof("create chat tab url %s", createChatTabUrl)
	req, err := http.NewRequest("POST", createChatTabUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	t, err := getToken()
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+t.TAToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}
	defer resp.Body.Close()

	res := &model.FeiShuResponse{}
	body, _ := io.ReadAll(resp.Body)
	logrus.Infof("response Body: %s", string(body))
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}
	return nil
}
