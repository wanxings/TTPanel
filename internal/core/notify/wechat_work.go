package notify

import (
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeChatWorkConfig struct {
	BotName string `json:"bot_name"`
	BotUrl  string `json:"bot_url"`
}

type WeChatWorkNotifier struct {
	config WeChatWorkConfig
}

func (n *WeChatWorkNotifier) GetConfig() interface{} {
	return struct {
		WeChatWorkConfig WeChatWorkConfig `json:"wechat_work_config"`
	}{
		WeChatWorkConfig: n.config,
	}
}

func NewWeChatNotifier(config WeChatWorkConfig) *WeChatWorkNotifier {
	return &WeChatWorkNotifier{config: config}
}

func (n *WeChatWorkNotifier) Send(level, title, content string) error {
	switch level {
	case constant.NotifyLevelInfo:
		title = fmt.Sprintf("[%s] %s", helper.Message("notify.Info"), title)
	case constant.NotifyLevelWarning:
		title = fmt.Sprintf("[%s] %s", helper.Message("notify.Warning"), title)
	case constant.NotifyLevelSuccess:
		title = fmt.Sprintf("[%s] %s", helper.Message("notify.Success"), title)
	case constant.NotifyLevelDebug:
		title = fmt.Sprintf("[%s] %s", helper.Message("notify.Debug"), title)
	default:
		title = fmt.Sprintf("[%s] %s", helper.Message("notify.Info"), title)
	}
	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("%s: %s", title, content),
		},
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(n.config.BotUrl, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return nil
}
