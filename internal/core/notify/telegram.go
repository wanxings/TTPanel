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

type TelegramConfig struct {
	BotName  string `json:"bot_name"`
	BotToken string `json:"bot_token"`
	ToChatId string `json:"to_chat_id"`
}

type TelegramNotifier struct {
	config TelegramConfig
}

func (n *TelegramNotifier) GetConfig() interface{} {
	return struct {
		TelegramConfig TelegramConfig `json:"telegram_config"`
	}{
		TelegramConfig: n.config,
	}
}

func NewTelegramNotifier(config TelegramConfig) *TelegramNotifier {
	return &TelegramNotifier{config: config}
}

func (n *TelegramNotifier) Send(level, title, content string) error {
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
		"chat_id": n.config.ToChatId, //
		"text":    fmt.Sprintf("【%s】: %s", title, content),
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	sendMessageURL := "https://api.telegram.org/bot" + n.config.BotToken + "/sendMessage"
	resp, err := http.Post(sendMessageURL, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return nil
}
