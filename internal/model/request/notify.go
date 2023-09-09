package request

import "TTPanel/internal/core/notify"

type AddNotifyChannelR struct {
	Category    int    `json:"category" form:"category" binding:"required"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	NotifyChannelConfigR
}

type NotifyChannelConfigR struct {
	EmailConfig      notify.EmailConfig      `json:"email_config" form:"email_config"`
	DingTalkConfig   notify.DingTalkConfig   `json:"ding_talk_config" form:"ding_talk_config"`
	WeChatWorkConfig notify.WeChatWorkConfig `json:"wechat_work_config" form:"wechat_work_config"`
	TelegramConfig   notify.TelegramConfig   `json:"telegram_config" form:"telegram_config"`
	AliSmsConfig     notify.AliSmsConfig     `json:"ali_sms_config" form:"ali_sms_config"`
}

//type EmailConfigR struct {
//	EmailFrom     string `json:"email_from" form:"email_from"`
//	EmailPassword string `json:"email_password" form:"email_password"`
//	EmailTo       string `json:"email_to" form:"email_to"`
//	EmailSmtpHost string `json:"email_smtp_host" form:"email_smtp_host"`
//	EmailSmtpPort int    `json:"email_smtp_port" form:"email_smtp_port"`
//}
//type DingTalkBotConfigR struct {
//	DingTalkBotName   string `json:"dingTalk_bot_name" form:"ding_bot_name"`
//	DingTalkBotUrl    string `json:"dingTalk_bot_url" form:"ding_bot_url"`
//	DingTalkBotSecret string `json:"dingTalk_bot_secret" form:"dingTalk_bot_secret"`
//}
//type WechatWorkBotConfigR struct {
//	WechatWorkBotName string `json:"wechatWork_bot_name" form:"wechatWork_bot_name"`
//	WechatWorkBotUrl  string `json:"wechatWork_bot_url" form:"wechatWork_bot_url"`
//}
//type TelegramBotConfigR struct {
//	TelegramBotName   string `json:"telegram_bot_name" form:"telegram_bot_name"`
//	TelegramBotToken  string `json:"telegram_bot_token" form:"telegram_bot_token"`
//	TelegramBotChatId string `json:"telegram_bot_chat_id" form:"telegram_bot_chat_id"`
//}
//type AliSmsConfigR struct {
//	AliSmsAccessKeyId     string `json:"aliSms_access_key_id" form:"aliSms_access_key_id"`
//	AliSmsAccessKeySecret string `json:"aliSms_access_key_secret" form:"aliSms_access_key_secret"`
//	AliSmsSignName        string `json:"aliSms_sign_name" form:"aliSms_sign_name"`
//	AliSmsTemplateCode    string `json:"aliSms_template_code" form:"aliSms_template_code"`
//	AliSmsTo              string `json:"aliSms_to" form:"aliSms_to"`
//}

type EditNotifyChannelR struct {
	ID int64 `json:"id" form:"id" binding:"required"`
	AddNotifyChannelR
}

type TestNotifyChannelR struct {
	Category int `json:"category" form:"category" binding:"required"`
	NotifyChannelConfigR
}

type NotifyChannelListR struct {
	Query    string `json:"query" form:"query"`
	Category int    `json:"category" form:"category"`
	Limit    int    `json:"limit" form:"limit" binding:"required"`
	Page     int    `json:"page" form:"page" binding:"required"`
}

type DeleteNotifyChannelR struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}
