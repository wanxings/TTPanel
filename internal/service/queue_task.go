package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/pkg/util"
	"fmt"
	"os"
)

type QueueTaskService struct{}

// 0未執行1執行中2執行完成3執行錯誤
//var (
//	taskQueueWg   sync.WaitGroup
//	quitTaskQueue chan string
//)

// CreateTaskQueue 添加任務至任務隊列
func (t *QueueTaskService) CreateTaskQueue(taskQueue *model.QueueTask) error {
	_, err := taskQueue.Create(global.PanelDB)
	if err != nil {
		global.Log.Errorf("添加任務失敗->createTaskQueue()->ds.CreateTaskQueue()  Error:%s", err)
		return err
	}
	return nil
}

// RunningCount 运行中的任务数量
func (t *QueueTaskService) RunningCount() (int64, error) {
	whereT := model.ConditionsT{
		"status IN (?)": []int{constant.QueueTaskStatusProcessing, constant.QueueTaskStatusWait},
	}
	return (&model.QueueTask{}).Count(global.PanelDB, whereT)
}

// TaskList 任务列表
func (t *QueueTaskService) TaskList(status, offset, limit int) ([]*model.QueueTask, int64, error) {
	list, total, err := (&model.QueueTask{}).List(global.PanelDB, &model.ConditionsT{"status = ?": status, "ORDER": "create_time DESC"}, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// GetTaskQueueList 获取任务列表
func (t *QueueTaskService) GetTaskQueueList() ([]*model.QueueTask, error) {
	// 标记上次未执行成功的任务
	err := (&model.QueueTask{}).UpdateOne(global.PanelDB, "status", constant.QueueTaskStatusWait, &model.ConditionsT{
		"status": constant.QueueTaskStatusProcessing,
	})
	if err != nil {
		global.Log.Errorf("GetTaskQueueList->task.QueueTask.UpdateOne  Error:%s", err)
		return nil, err
	}
	// 查询未执行的任务
	taskQueueList, _, err := (&model.QueueTask{}).List(global.PanelDB, &model.ConditionsT{
		"status": constant.QueueTaskStatusWait,
	}, 0, 0)
	if err != nil {
		global.Log.Errorf("查询未执行的任务  Error:%s", err)
		return nil, err
	}
	return taskQueueList, nil
}

// GetQueueTaskByID 根据ID获取任务
func (t *QueueTaskService) GetQueueTaskByID(id int64) (*model.QueueTask, error) {
	return (&model.QueueTask{ID: id}).Get(global.PanelDB)
}

// DelTask 删除任务
func (t *QueueTaskService) DelTask(taskQueue *model.QueueTask) error {
	switch taskQueue.Status {
	case constant.QueueTaskStatusWait:
		//直接删除
		return taskQueue.Delete(global.PanelDB)
	case constant.QueueTaskStatusSuccess:
		//直接删除
		err := taskQueue.Delete(global.PanelDB)
		if err != nil {
			return err
		}
		//清理日志
		_ = os.Remove(fmt.Sprintf("%s/data/logs/queue_task/%d_%d.log", global.Config.System.PanelPath, taskQueue.ID, taskQueue.CreateTime))
		return nil

	case constant.QueueTaskStatusProcessing:
		//直接删除
		err := taskQueue.Delete(global.PanelDB)
		if err != nil {
			return err
		}
		//终止相关进程
		cmdStr := fmt.Sprintf("target_command=\"%s\" \n", taskQueue.ExecStr)
		cmdStr += `
pid_list=$(ps aux | grep "$target_command" | grep -v grep | awk '{print $2}')
for pid in $pid_list; do
    kill -9 $pid
done
`
		_, err = util.ExecShellScript(cmdStr)
		if err != nil {
			global.Log.Debugf("终止相关进程失败->DelTask()->util.ExecShellScript()  Error:%s", err)
		}
		//清理日志
		_ = os.Remove(fmt.Sprintf("%s/data/logs/queue_task/%d_%d.log", global.Config.System.PanelPath, taskQueue.ID, taskQueue.CreateTime))
		return err
	default:
		return nil
	}
}

// ClearTask 清理任务
func (t *QueueTaskService) ClearTask() error {
	//查询任务
	taskQueueList, _, err := (&model.QueueTask{}).List(global.PanelDB, &model.ConditionsT{
		"status": constant.QueueTaskStatusSuccess,
	}, 0, 0)
	if err != nil {
		global.Log.Errorf("查询未执行的任务  Error:%s", err)
		return err
	}
	for _, taskQueue := range taskQueueList {
		//直接删除
		err = taskQueue.Delete(global.PanelDB)
		if err != nil {
			return err
		}
		//清理日志
		_ = os.Remove(fmt.Sprintf("%s/data/logs/queue_task/%d_%d.log", global.Config.System.PanelPath, taskQueue.ID, taskQueue.CreateTime))
	}
	return nil
}
