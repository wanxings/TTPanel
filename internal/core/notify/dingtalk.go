package notify

import (
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DingTalkConfig struct {
	BotName string `json:"bot_name"`
	BotUrl  string `json:"bot_url"`
	Secret  string `json:"secret"`
}

type DingTalkNotifier struct {
	config DingTalkConfig
}

func (n *DingTalkNotifier) GetConfig() interface{} {
	return struct {
		DingTalkConfig DingTalkConfig `json:"dingTalk_config"`
	}{
		DingTalkConfig: n.config,
	}
}

func NewDingTalkNotifier(config DingTalkConfig) *DingTalkNotifier {
	return &DingTalkNotifier{config: config}
}

func (n *DingTalkNotifier) Send(level, title, content string) error {
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
			"content": fmt.Sprintf("【%s】: %s", title, content),
		},
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	timestamp := time.Now().UnixNano() / 1e6
	nonce := "123456" // 随机字符串，可以自定义或者使用随机生成的字符串

	// 将消息体和时间戳拼接成一个字符串，用于加密
	messageString := fmt.Sprintf("%d\n%s\n%s", timestamp, nonce, string(messageJSON))

	// 对消息体进行加密
	encryptedMessage, err := encryptAES(messageString, []byte(n.config.Secret), []byte(nonce))
	if err != nil {
		return err
	}

	// 对加密后的消息进行Base64编码
	base64Message := base64.StdEncoding.EncodeToString(encryptedMessage)

	// 构造请求参数
	requestURL := fmt.Sprintf("%s&timestamp=%d&nonce=%s&sign=%s", n.config.BotUrl, timestamp, nonce, base64Message)

	// 发送请求
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return nil
}

// 使用AES-CBC算法对消息进行加密
func encryptAES(message string, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 补全消息长度，使其长度为16的倍数
	messageBytes := usePKCS7Padding([]byte(message), block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(messageBytes))
	blockMode.CryptBlocks(encrypted, messageBytes)

	return encrypted, nil
}

// usePKCS7Padding 使用PKCS7算法对消息进行填充
func usePKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
