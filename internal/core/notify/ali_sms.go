package notify

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type AliSmsConfig struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SignName        string `json:"sign_name"`
	TemplateCode    string `json:"template_code"`
	To              string `json:"to"`
}

type AliSmsNotifier struct {
	config AliSmsConfig
}

func (n *AliSmsNotifier) GetConfig() interface{} {
	return struct {
		AliSmsConfig AliSmsConfig `json:"ali_sms_config"`
	}{
		AliSmsConfig: n.config,
	}
}

func NewAliSmsNotifier(config AliSmsConfig) *AliSmsNotifier {
	return &AliSmsNotifier{config: config}
}

func (n *AliSmsNotifier) Send(level, title, content string) error {
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

	client, _err := n.CreateClient()
	if _err != nil {
		return _err
	}

	var templateParam struct {
		Msg string `json:"msg"`
	}
	templateParam.Msg = fmt.Sprintf("%s-%s", title, content)
	// 将结构体转换为JSON字符串
	templateParamStr, err := json.Marshal(templateParam)
	if err != nil {
		global.Log.Error(fmt.Sprintf("alisms->Send->json.Marshal failed, err:%v", err))
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(n.config.SignName),
		TemplateCode:  tea.String(n.config.TemplateCode),
		PhoneNumbers:  tea.String(n.config.To),
		TemplateParam: tea.String(string(templateParamStr)),
	}
	runtime := &util.RuntimeOptions{}
	sendSmsResponse, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if _err != nil {
		global.Log.Errorf("alisms->Send->client.SendSmsWithOptions failed, err:%v", _err)
		return _err
	}
	if sendSmsResponse.Body.Code != tea.String("OK") {
		global.Log.Errorf("alisms->Send->sendSmsResponse.Body.Code:%v,Message:%v", sendSmsResponse.Body.Code, sendSmsResponse.Body.Message)
		return _err
	}

	return nil
}
func (n *AliSmsNotifier) CreateClient() (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(n.config.AccessKeyId),
		AccessKeySecret: tea.String(n.config.AccessKeySecret),
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
