package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/pkg/util"
	"fmt"
)

// ReadExtensionsInfo 读取扩展信息
func ReadExtensionsInfo(name string, info interface{}) error {
	infoBody, err := util.ReadFileStringBody(global.Config.System.PanelPath + "/data/extensions/" + name + "/info.json")
	if err != nil {
		return err
	}
	err = util.JsonStrToStruct(infoBody, info)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// GetExtensionsPath 取扩展路径
func GetExtensionsPath(ExtensionName string) string {
	return global.Config.System.PanelPath + "/data/extensions/" + ExtensionName
}

// GetExtensionsShellPath 取扩展shell路径
func GetExtensionsShellPath(ExtensionName string, shellName string) string {
	return fmt.Sprintf("%s/%s", GetExtensionsPath(ExtensionName), shellName)
}

// AddTaskQueue 添加任務至任務隊列
func AddTaskQueue(Name, ExecStr string) error {
	_, err := (&model.QueueTask{
		Name:    Name,
		Type:    1,
		Status:  constant.QueueTaskStatusWait,
		ExecStr: ExecStr,
	}).Create(global.PanelDB)
	if err != nil {
		global.Log.Errorf("添加任務失敗->createTaskQueue()->ds.CreateTaskQueue()  Error:%s", err)
		return err
	}
	return nil
}

// CheckTaskQueueExists 检查任务队列是否存在
func CheckTaskQueueExists(taskName string) (bool, error) {
	_, total, err := (&model.QueueTask{}).List(global.PanelDB, &model.ConditionsT{
		"FIXED": "status = " + fmt.Sprintf("%d", constant.QueueTaskStatusProcessing) + " OR status = " + fmt.Sprintf("%d", constant.QueueTaskStatusWait),
		"name":  taskName,
	}, 0, 0)
	if err != nil {
		return false, err
	}
	if total > 0 {
		return true, nil
	}
	return false, nil
}
