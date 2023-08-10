package dao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lambda-lark-bot/config"
	"lambda-lark-bot/model"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	sendMsgUrl       = "https://open.feishu.cn/open-apis/message/v4/send/"
	downloadUrl      = "https://open.feishu.cn/open-apis/im/v1/messages/%s/resources/%s?type=%s"
	createChannelUrl = "https://open.feishu.cn/open-apis/chat/v4/create/"
	tokenUrl         = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal/"
	createChatTabUrl = "https://open.feishu.cn/open-apis/im/v1/chats/%s/chat_tabs"
)

func CreateChannel(userIDs []string, name string) (channelID string, err error) {
	m := &model.CreateChannelReq{
		UserIds: userIDs,
		Name:    name,
	}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", createChannelUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	t, err := getToken()
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+t.TAToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("%v", err)
		return "", err
	}
	defer resp.Body.Close()

	res := &model.FeiShuResponse{}
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Infof("response Body: %s", string(body))
	err = json.Unmarshal(body, res)
	if err != nil {
		return "", err
	}
	return res.Data.ChatID, nil
}

func SendFeishuMsg(feishuMsg *model.FeiShuMsg) ([]byte, error) {
	jsonStr, err := json.Marshal(feishuMsg)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", sendMsgUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	t, err := getToken()
	if err != nil {
		logrus.Errorf("failed to get token %+v", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+t.TAToken)
	req.Header.Set("Content-Type", "application/json")

	logrus.Infof("send msg req %v", req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Infof("send msg response Body: %s", string(body))
	return body, nil
}

func SendMsgByType(chatId, userID, msg, msgType string) ([]byte, error) {
	fsmsg := &model.FeiShuMsg{
		MsgType: msgType,
		Content: &model.Content{
			Text: msg,
		},
		UserId: userID,
		ChatId: chatId,
	}

	return SendFeishuMsg(fsmsg)
}

//SendMsg chatId group ID
//SendMsg userID userID
func SendMsg(chatId, userID, msg string) error {
	_, err := SendMsgByType(chatId, userID, msg, "text")
	return err
}

func SendCardMsg(msgCard *model.FeiShuMsg, caze *Case) (*model.CardRsp, error) {
	tmpElement := msgCard.Card.Elements[0]
	model.BuildCardWithTitle(&msgCard.Card, caze.Title)
	model.BuildCardWithContent(&msgCard.Card, caze.Content)

	respBody, err := SendFeishuMsg(msgCard)
	if err != nil {
		logrus.Errorf("Failed to send card msg, %v", err)
		return nil, err
	}
	rsp := &model.CardRsp{}
	err = json.Unmarshal(respBody, rsp)
	if err != nil {
		logrus.Errorf("Failed to unmarshal, %v", err)
		return nil, err
	}
	msgCard.Card.Elements[0] = tmpElement
	return rsp, nil
}

func SendErrCardMsg(chatId, userID string, e error) error {
	config.Conf.ErrCardTemplate.Card.Elements[0].Content = e.Error()
	config.Conf.ErrCardTemplate.ChatId = chatId
	rsp, err := SendFeishuMsg(config.Conf.ErrCardTemplate)
	logrus.Infof("Send err card rsp %v", string(rsp))
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

	body, _ := ioutil.ReadAll(resp.Body)
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
		body, _ := ioutil.ReadAll(resp.Body)
		logrus.Infof("bad status: %s", string(body))
		return nil, fmt.Errorf("bad status: %v", resp)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
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
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Infof("response Body: %s", string(body))
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}
	return nil
}
