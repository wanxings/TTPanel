package notify

import (
	"TTPanel/internal/helper/constant"
	"fmt"
)

type Notify interface {
	Send(level string, title, content string) error
	GetConfig() interface{}
}

func New(notifyType int, config interface{}) (Notify, error) {
	//var notifier Notifier
	var err error

	switch notifyType {
	case constant.NotifyCategoryByEmail:
		emailConfig, ok := config.(EmailConfig)
		if !ok {
			err = fmt.Errorf("invalid email config")
			break
		}
		return NewEmailNotifier(emailConfig), nil
	case constant.NotifyCategoryByDingTalk:
		dingTalkConfig, ok := config.(DingTalkConfig)
		if !ok {
			err = fmt.Errorf("invalid dingtalk config")
			break
		}
		return NewDingTalkNotifier(dingTalkConfig), nil
	case constant.NotifyCategoryByWeChatWork:
		weChatConfig, ok := config.(WeChatWorkConfig)
		if !ok {
			err = fmt.Errorf("invalid wechat config")
			break
		}
		return NewWeChatNotifier(weChatConfig), nil
	case constant.NotifyCategoryByAliSms:
		aliSmsConfig, ok := config.(AliSmsConfig)
		if !ok {
			err = fmt.Errorf("invalid aliSms config")
			break
		}
		return NewAliSmsNotifier(aliSmsConfig), nil
	case constant.NotifyCategoryByTelegram:
		telegramConfig, ok := config.(TelegramConfig)
		if !ok {
			err = fmt.Errorf("invalid telegram config")
			break
		}
		return NewTelegramNotifier(telegramConfig), nil
	default:
		err = fmt.Errorf("invalid notify type")
	}

	return nil, err
}
