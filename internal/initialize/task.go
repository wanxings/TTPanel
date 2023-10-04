package initialize

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var (
	taskQueueWg   sync.WaitGroup
	quitTaskQueue chan string
)

func TaskInit() {
	quitTaskQueue = make(chan string) //退出通道
	taskQueueWg.Add(1)
	go startTaskQueue(quitTaskQueue)
	_, _ = fmt.Fprintf(color.Output, "Queue task....   %s\n",
		color.GreenString("done"),
	)
}

func startTaskQueue(quit <-chan string) {
	defer taskQueueWg.Done()
	for {
		select {
		case <-quit:
			_, _ = fmt.Fprintf(color.Output, "Queue task quit....   %s\n",
				color.GreenString("done"),
			)
			return //必须return，否则goroutine不会结束
		default:
			//taskQueueList, err := service.GroupApp.QueueTaskServiceApp.GetTaskQueueList()
			// 标记上次未执行成功的任务
			err := (&model.QueueTask{}).UpdateOne(global.PanelDB, "status", constant.QueueTaskStatusWait, &model.ConditionsT{
				"status": constant.QueueTaskStatusProcessing,
			})
			if err != nil {
				global.Log.Errorf("GetTaskQueueList->task.QueueTask.UpdateOne  Error:%s", err)
				return
			}
			// 查询未执行的任务
			taskQueueList, _, err := (&model.QueueTask{}).List(global.PanelDB, &model.ConditionsT{
				"status": constant.QueueTaskStatusWait,
				"ORDER":  "create_time ASC",
			}, 0, 0)
			if err != nil {
				global.Log.Errorf("查询未执行的任务  Error:%s", err)
				return
			}

			if err != nil {
				global.Log.Errorf("GetTaskQueueList  Error:%s", err)
				return
			}

			for _, taskQueue := range taskQueueList {
				global.Log.Debugln("QueueTask:", taskQueue)
				taskQueue.StartTime = time.Now().Unix()
				taskQueue.Status = constant.QueueTaskStatusProcessing

				err := taskQueue.Update(global.PanelDB)
				if err != nil {
					global.Log.Errorf("StartTaskQueue->UpdateTaskQueueBatch  Error:%s", err)
					continue
				}
				logPath := fmt.Sprintf("%s/queue_task/%d_%d.log", global.Config.Logger.RootPath, taskQueue.ID, taskQueue.CreateTime)
				_ = os.MkdirAll(filepath.Dir(logPath), os.ModePerm)
				cmdStr := fmt.Sprintf("/bin/bash -c '%s' > %s 2>&1", taskQueue.ExecStr, logPath)
				global.Log.Debugf("QueueTaskCmdStr:%s", cmdStr)
				cmd := exec.Command("/bin/bash", "-c", cmdStr)
				err = cmd.Run()
				if err != nil {
					global.Log.Errorf("StartTaskQueue->cmd.Run  Error:%s", err)
				}
				taskQueue.EndTime = time.Now().Unix()
				taskQueue.Status = constant.QueueTaskStatusSuccess
				err = taskQueue.Update(global.PanelDB)
				if err != nil {
					global.Log.Errorf("StartTaskQueue->UpdateTaskQueueBatch  Error:%s", err)
					continue
				}
				global.Log.Debugf("QueueTask %s done,time cost %s", taskQueue.Name, util.ResolveTime(taskQueue.EndTime, taskQueue.StartTime))
			}
			time.Sleep(5 * time.Second)
		}
	}

}
