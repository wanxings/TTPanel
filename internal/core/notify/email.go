package notify

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
)

type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	From     string `json:"from"`
	To       string `json:"to"`
	Password string `json:"password"`
}

type EmailNotifier struct {
	config EmailConfig
}

func (n *EmailNotifier) GetConfig() interface{} {
	return struct {
		EmailConfig EmailConfig `json:"email_config"`
	}{
		EmailConfig: n.config,
	}
}

func NewEmailNotifier(config EmailConfig) *EmailNotifier {
	return &EmailNotifier{config: config}
}

func (n *EmailNotifier) Send(level, title, content string) error {
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
	// 读取本地html模板文件
	global.Log.Debugf("Send->n.config:%v \n", n.config)
	tmpl, err := template.ParseFiles(global.Config.System.PanelPath + "/data/notify/template/email.html")
	if err != nil {
		return err
	}

	// 定义模板数据
	data := struct {
		Title     string
		Content   string
		PanelName string
	}{
		Title:     title,
		Content:   content,
		PanelName: global.Config.System.PanelName,
	}

	// 渲染模板并将结果写入缓冲区
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, data); err != nil {
		return err
	}
	// 将缓冲区中的渲染结果作为邮件正文
	body := tpl.String()

	m := gomail.NewMessage()

	//设置发件人
	m.SetHeader("From", n.config.From)

	//设置收件用户
	m.SetHeader("To", n.config.To)

	//设置邮件主题
	m.SetHeader("Subject", title)

	//设置邮件正文
	m.SetBody("text/html", body)
	d := gomail.NewDialer(n.config.Host, n.config.Port, n.config.From, n.config.Password)

	return d.DialAndSend(m)
}
