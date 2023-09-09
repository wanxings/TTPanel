package service

import (
	notifyCore "TTPanel/internal/core/notify"
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
)

type NotifyService struct {
}

// AddNotifyChannel 添加通知通道
func (s *NotifyService) AddNotifyChannel(param *request.AddNotifyChannelR) error {
	notify, err := s.NewNotifyChannelCore(param.Category, &param.NotifyChannelConfigR)
	if err != nil {
		return err
	}
	//将结构体转换为json字符串
	notifyConfigStr, err := util.StructToJsonStr(notify.GetConfig())
	if err != nil {
		return err
	}

	//构造通知数据
	notifyData := &model.NotifyChannel{
		Category:    param.Category,
		Name:        param.Name,
		Description: param.Description,
		Config:      notifyConfigStr,
	}
	//插入通知数据
	_, err = notifyData.Create(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// EditNotifyChannel 编辑通知通道
func (s *NotifyService) EditNotifyChannel(param *request.EditNotifyChannelR) error {
	//查询通知数据
	notifyOldData, err := (&model.NotifyChannel{ID: param.ID}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if notifyOldData.ID == 0 {
		return errors.New("not found NotifyChannel")
	}

	notify, err := s.NewNotifyChannelCore(notifyOldData.Category, &param.NotifyChannelConfigR)
	if err != nil {
		return err
	}

	notifyOldData.Name = param.Name
	notifyOldData.Description = param.Description
	notifyOldData.Category = param.Category
	notifyOldData.Config, _ = util.StructToJsonStr(notify.GetConfig())

	//更新通知数据
	err = notifyOldData.Update(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// NotifyChannelList 获取通知通道列表
func (s *NotifyService) NotifyChannelList(query string, category int, offset, limit int) ([]*model.NotifyChannel, int64, error) {
	whereT := model.ConditionsT{"ORDER": "create_time DESC"}
	whereOrT := model.ConditionsT{}
	if !util.StrIsEmpty(query) {
		query = "%" + query + "%"
		whereT["name LIKE ?"] = query
		whereOrT["description LIKE ?"] = query
	}
	if category > 0 {
		whereT["category"] = category
	}
	return (&model.NotifyChannel{}).List(global.PanelDB, &whereT, &whereOrT, offset, limit)
}

// TestNotifyChannel 测试通知通道
func (s *NotifyService) TestNotifyChannel(param *request.TestNotifyChannelR) error {
	//发送通知
	notify, err := s.NewNotifyChannelCore(param.Category, &param.NotifyChannelConfigR)
	if err != nil {
		return err
	}
	err = notify.Send(constant.NotifyLevelDebug, "Test Title", "Test Content")
	if err != nil {
		return err
	}
	return nil
}

func (s *NotifyService) SendNotify(notifyID int64, level string, title, content string) error {
	//查询通知配置
	notifyGet, err := (&model.NotifyChannel{ID: notifyID}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if notifyGet.ID == 0 {
		fmt.Println(fmt.Errorf("Not Found Notify By ID: %d \n", notifyID))
		return err
	}
	var config request.NotifyChannelConfigR
	err = util.JsonStrToStruct(notifyGet.Config, &config)
	if err != nil {
		return err
	}
	notify, err := s.NewNotifyChannelCore(notifyGet.Category, &config)
	if err != nil {
		return err
	}
	err = notify.Send(level, title, content)
	if err != nil {
		return err
	}
	return nil
}

// NewNotifyChannelCore 创建通知通道核心
func (s *NotifyService) NewNotifyChannelCore(category int, param *request.NotifyChannelConfigR) (notifyCore.Notify, error) {
	var notifyConfig interface{}
	switch category {
	case constant.NotifyCategoryByEmail:
		notifyConfig = param.EmailConfig
	case constant.NotifyCategoryByWeChatWork:
		notifyConfig = param.WeChatWorkConfig
	case constant.NotifyCategoryByDingTalk:
		notifyConfig = param.DingTalkConfig
	case constant.NotifyCategoryByAliSms:
		notifyConfig = param.AliSmsConfig
	case constant.NotifyCategoryByTelegram:
		notifyConfig = param.TelegramConfig
	default:
		return nil, errors.New(fmt.Sprintf("unknown notify category: %d", category))
	}
	return notifyCore.New(category, notifyConfig)
}

// DeleteNotifyChannel 删除通知通道
func (s *NotifyService) DeleteNotifyChannel(id int64) error {
	//查询通知数据
	notifyData, err := (&model.NotifyChannel{ID: id}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if notifyData.ID == 0 {
		return errors.New("not found NotifyChannel")
	}
	//删除通知数据
	err = notifyData.Delete(global.PanelDB, &model.ConditionsT{})
	if err != nil {
		return err
	}
	return nil
}
