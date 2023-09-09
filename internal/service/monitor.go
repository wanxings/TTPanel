package service

import (
	"TTPanel/internal/conf"
	"TTPanel/internal/core/system"
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/fatih/color"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var (
	monitorWg         sync.WaitGroup
	monitorGisRunning atomic.Value
	quitMonitorG      chan string
)
var eventMutex = &sync.Mutex{}
var EventConfig = &request.MonitorEventConfig{}

type MonitorService struct{}

func MonitorInit() {
	if monitorGisRunning.Load() == true {
		// 已在运行中,直接返回
		return
	}
	if global.Config.Monitor.Status {
		//监控配置
		configBody, err := util.ReadFileStringBody(global.Config.System.PanelPath + "/data/monitor/config.json")
		if err != nil {
			panic(err)
			return
		}
		err = util.JsonStrToStruct(configBody, EventConfig)
		if err != nil {
			panic(err)
			return
		}

		// 运行监控
		monitorGisRunning.Store(true)
		quitMonitorG = make(chan string)
		monitorWg.Add(1)
		go startMonitor(quitMonitorG)
		_, _ = fmt.Fprintf(color.Output, "System monitoring....   %s\n",
			color.GreenString("done"),
		)
	}
}

func StopMonitor() {
	quitMonitorG <- "quit"
	monitorWg.Wait()
}

func startMonitor(quit <-chan string) {
	defer monitorWg.Done()
	for {
		select {
		case <-quit:
			_, _ = fmt.Fprintf(color.Output, "System monitoring quit....   %s\n",
				color.GreenString("done"),
			)
			monitorGisRunning.Store(false)
			return //必须return，否则goroutine不会结束
		default:
			time.Sleep(10 * time.Second)
			cpuInfo, _ := system.GetCPU()
			ioStat := system.GetIoStat()
			memory := system.GetMemory()
			NetworkData := system.GetNetWork()
			avgStat := system.GetLoad()
			disk := system.GetDisk()

			//保存MonitorNetwork
			_, err := (&model.MonitorNetwork{
				Up:          NetworkData["all"].Up,
				Down:        NetworkData["all"].Down,
				TotalUp:     NetworkData["all"].UpTotal,
				TotalDown:   NetworkData["all"].DownTotal,
				DownPackets: NetworkData["all"].DownPackets,
				UpPackets:   NetworkData["all"].UpPackets,
			}).Create(global.PanelDB)
			if err != nil {
				global.Log.Errorf("startMonitor->model.MonitorNetwork.Create Error:%v \n", err)
				return
			}
			delTime := time.Now().Add(-time.Duration(global.Config.Monitor.SaveDay) * 24 * time.Hour).Unix()
			//清理过期的MonitorNetwork
			err = (&model.MonitorNetwork{}).BatchDelete(global.PanelDB, &model.ConditionsT{"create_time < ?": delTime})
			if err != nil {
				global.Log.Errorf("startMonitor->systemModel.MonitorNetwork.BatchDelete Error:%v \n", err)
			}

			//保存MonitorMetrics
			cpuUsage := 0.0
			for _, usage := range cpuInfo.Percents {
				cpuUsage += usage
			}

			metrics, err := (&model.MonitorMetrics{
				CpuUsage:            float64(int((cpuUsage/float64(len(cpuInfo.Percents)))*100)) / 100,
				MemUsage:            float64(int(memory.Percent*100)) / 100,
				LoadAvg1m:           avgStat.Load1,
				LoadAvg5m:           avgStat.Load5,
				LoadAvg15m:          avgStat.Load15,
				ResourceUtilization: float64(int((avgStat.Load1/(float64(cpuInfo.CoresCount*2)*0.75)*100)*100)) / 100,
			}).Create(global.PanelDB)
			if err != nil {
				global.Log.Errorf("startMonitor->systemModel.MonitorMetrics.Create Error:%v \n", err)
				return
			}
			//清理过期的MonitorMetrics
			err = (&model.MonitorMetrics{}).BatchDelete(global.PanelDB, &model.ConditionsT{"create_time < ?": delTime})
			if err != nil {
				global.Log.Errorf("startMonitor->systemModel.MonitorMetrics.BatchDelete Error:%v \n", err)
			}

			//保存MonitorIo
			_, err = (&model.MonitorIo{
				ReadCount:  ioStat["all"].ReadCount,
				WriteCount: ioStat["all"].WriteCount,
				ReadBytes:  ioStat["all"].ReadBytes,
				WriteBytes: ioStat["all"].WriteBytes,
				ReadTime:   ioStat["all"].ReadTime,
				WriteTime:  ioStat["all"].WriteTime,
			}).Create(global.PanelDB)
			if err != nil {
				global.Log.Errorf("startMonitor->systemModel.MonitorIo.Create Error:%v \n", err)
				return
			}
			//清理过期的MonitorIo
			err = (&model.MonitorIo{}).BatchDelete(global.PanelDB, &model.ConditionsT{"create_time < ?": delTime})
			if err != nil {
				global.Log.Errorf("startMonitor->systemModel.MonitorIo.BatchDelete Error:%v \n", err)
			}

			//监控事件处理
			monitorService := &MonitorService{}
			err = monitorService.CpuEvent(metrics.CpuUsage)
			if err != nil {
				global.Log.Errorf("startMonitor->CpuEvent Error:%v \n", err)
			}
			err = monitorService.MemEvent(metrics.MemUsage)
			if err != nil {
				global.Log.Errorf("startMonitor->MemEvent Error:%v \n", err)
			}
			err = monitorService.DiskSpaceEvent(disk)
			if err != nil {
				global.Log.Errorf("startMonitor->DiskSpaceEvent Error:%v \n", err)
			}
			err = monitorService.DiskInodeEvent(disk)
			if err != nil {
				global.Log.Errorf("startMonitor->DiskInodeEvent Error:%v \n", err)
			}
			//Todo: 服务监控,待优化
			serviceList := []string{constant.MonitorEventCategoryByServiceNginx, constant.MonitorEventCategoryByServiceMysql, constant.MonitorEventCategoryByServiceRedis, constant.MonitorEventCategoryByServiceDocker}
			for _, serviceName := range serviceList {
				err = monitorService.ServiceStatusEvent(serviceName)
				if err != nil {
					global.Log.Errorf("startMonitor->ServiceStatusEvent Error:%v \n", err)
				}
			}
		}
	}

}

// Base 获取基础信息
func (s *MonitorService) Base() (map[string]interface{}, error) {
	ResponseData := make(map[string]interface{})
	ResponseData["cpu"], ResponseData["cpu_times"] = system.GetCPU()
	ResponseData["disk"] = system.GetDisk()
	ResponseData["iostat"] = system.GetIoStat()
	ResponseData["mem"] = system.GetMemory()
	ResponseData["system"] = system.GetSystemInfo()
	ResponseData["network"] = system.GetNetWork()
	ResponseData["load"] = system.GetLoad()
	ResponseData["swap_mem"] = system.GetSwapMemory()
	return ResponseData, nil
}

// Logs 获取日志
func (s *MonitorService) Logs(sort string, startTime, endTime int64) (map[string]interface{}, error) {
	// all cpu mem io network
	backData := make(map[string]interface{})
	if sort == "metrics" {
		metricsList, err := (&model.MonitorMetrics{}).List(global.PanelDB, &model.ConditionsT{}, startTime, endTime, 0, 0)
		if err != nil {
			return nil, err
		}
		backData["metrics"] = metricsList
	} else if sort == "io" {
		ioList, err := (&model.MonitorIo{}).List(global.PanelDB, &model.ConditionsT{}, startTime, endTime, 0, 0)
		if err != nil {
			return nil, err
		}
		backData["io"] = ioList
	} else if sort == "network" {
		networkList, err := (&model.MonitorNetwork{}).List(global.PanelDB, &model.ConditionsT{}, startTime, endTime, 0, 0)
		if err != nil {
			return nil, err
		}
		backData["network"] = networkList
	} else {
		metricsList, err := (&model.MonitorMetrics{}).List(global.PanelDB, &model.ConditionsT{}, startTime, endTime, 0, 0)
		if err != nil {
			return nil, err
		}
		backData["metrics"] = metricsList
		ioList, err := (&model.MonitorIo{}).List(global.PanelDB, &model.ConditionsT{}, startTime, endTime, 0, 0)
		if err != nil {
			return nil, err
		}
		backData["io"] = ioList
		networkList, err := (&model.MonitorNetwork{}).List(global.PanelDB, &model.ConditionsT{}, startTime, endTime, 0, 0)
		if err != nil {
			return nil, err
		}
		backData["network"] = networkList
	}
	return backData, nil
}

// SaveConfig 保存配置
func (s *MonitorService) SaveConfig(config *request.MonitorConfigR) error {
	global.Vp.Set("monitor", conf.Monitor{
		Status:  config.Status,
		SaveDay: config.SaveDay,
	})

	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		global.Log.Errorf("SaveConfig->WriteConfig Error:%s \n", err.Error())
		os.Exit(1)
	}
	if config.Status {
		MonitorInit()
	} else {
		StopMonitor()
	}
	return nil
}

// ClearAllLogs 清空日志
func (s *MonitorService) ClearAllLogs() error {
	err := (&model.MonitorMetrics{}).DeleteAll(global.PanelDB)
	if err != nil {
		return err
	}
	err = (&model.MonitorIo{}).DeleteAll(global.PanelDB)
	if err != nil {
		return err
	}
	err = (&model.MonitorNetwork{}).DeleteAll(global.PanelDB)
	if err != nil {
		return err
	}
	return nil
}

// GetEventConfig 获取事件配置
func (s *MonitorService) GetEventConfig() *request.MonitorEventConfig {
	return EventConfig
}

// SaveEventConfig 保存事件配置
func (s *MonitorService) SaveEventConfig(config *request.MonitorEventConfig) error {
	eventMutex.Lock()
	defer eventMutex.Unlock()
	configBody, err := util.StructToJsonStr(config)
	if err != nil {
		return err
	}
	EventConfig = config
	return util.WriteFile(global.Config.System.PanelPath+"/data/monitor/config.json", []byte(configBody), 0644)
}

// EventList 事件列表
func (s *MonitorService) EventList(status, offset, limit int) ([]*model.MonitorEvent, int64, error) {
	return (&model.MonitorEvent{}).List(global.PanelDB, &model.ConditionsT{"status": status, "ORDER": "create_time DESC"}, offset, limit)
}

// BatchSetEventStatus 批量设置事件状态
func (s *MonitorService) BatchSetEventStatus(ids []int64, status int) error {
	for _, id := range ids {
		eventGet, err := (&model.MonitorEvent{ID: id}).Get(global.PanelDB)
		if err != nil {
			return err
		}
		eventGet.Status = status
		err = eventGet.Update(global.PanelDB)
		if err != nil {
			return err
		}
	}
	return nil
}

// CpuEvent cpu事件
func (s *MonitorService) CpuEvent(threshold float64) error {
	//如果通知状态为开启并且超过阈值并且设置了通知id
	if EventConfig.Cpu.Status && threshold > EventConfig.Cpu.Threshold && EventConfig.Cpu.NotifyId != 0 {
		//查询对应事件状态为未处理的数量
		count, err := (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"category": constant.MonitorEventCategoryByCPU, "status": 0, "FIXED": fmt.Sprintf("create_time > %d", util.GetDayTimestampZ())})
		if err != nil {
			return err
		}
		if EventConfig.Cpu.NotifyCount <= count {
			return nil
		}
		//获取当前进程cpu使用
		cpuTop, _ := util.ExecShell("ps aux --sort=-%cpu | head -n 6")
		//插入事件
		_, err = (&model.MonitorEvent{
			Category: constant.MonitorEventCategoryByCPU,
			Log:      fmt.Sprintf("CPU use: %v %% \r\n", threshold) + cpuTop,
			Status:   0,
		}).Create(global.PanelDB)
		if err != nil {
			return err
		}
		//通知
		err = (&NotifyService{}).SendNotify(
			EventConfig.Cpu.NotifyId,
			constant.NotifyLevelWarning,
			fmt.Sprintf("Panel[%v] CPU threshold", global.Config.System.PanelName),
			fmt.Sprintf("CPU use: %v%% ", threshold))
		if err != nil {
			return err
		}
	}
	return nil
}

// MemEvent Mem事件
func (s *MonitorService) MemEvent(threshold float64) error {
	//如果通知状态为开启并且超过阈值并且设置了通知id
	if EventConfig.Mem.Status && threshold > EventConfig.Mem.Threshold && EventConfig.Mem.NotifyId != 0 {
		//查询对应事件状态为未处理的数量
		count, err := (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"category": constant.MonitorEventCategoryByMemory, "status": 0, "FIXED": fmt.Sprintf("create_time > %d", util.GetDayTimestampZ())})
		if err != nil {
			return err
		}
		if EventConfig.Mem.NotifyCount <= count {
			return nil
		}
		//获取当前进程cpu使用
		memTop, _ := util.ExecShell("ps aux --sort=-%mem | head -n 6")
		//插入事件
		_, err = (&model.MonitorEvent{
			Category: constant.MonitorEventCategoryByMemory,
			Log:      memTop,
			Status:   0,
		}).Create(global.PanelDB)
		if err != nil {
			return err
		}
		//通知
		err = (&NotifyService{}).SendNotify(
			EventConfig.Mem.NotifyId,
			constant.NotifyLevelWarning,
			fmt.Sprintf("Panel[%v] Memory threshold", global.Config.System.PanelName),
			fmt.Sprintf("Memory use: %v%% ", threshold))
		if err != nil {
			return err
		}
	}
	return nil
}

// DiskSpaceEvent DiskSpace事件
func (s *MonitorService) DiskSpaceEvent(diskList *[]response.Disk) error {
	if !EventConfig.DiskSpace.Status || EventConfig.DiskSpace.NotifyId == 0 {
		return nil
	}
	for _, mountPoint := range EventConfig.DiskSpace.MountPoint {
		for _, disk := range *diskList {
			if disk.MountPoint == mountPoint && EventConfig.DiskSpace.Threshold > disk.SpaceUsage.Percent {
				//查询对应事件状态为未处理的数量
				count, err := (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"category": constant.MonitorEventCategoryByDiskSpace, "status": 0, "FIXED": fmt.Sprintf("create_time > %d", util.GetDayTimestampZ())})
				if err != nil {
					return err
				}
				if EventConfig.DiskSpace.NotifyCount <= count {
					return nil
				}
				//插入事件
				_, err = (&model.MonitorEvent{
					Category: constant.MonitorEventCategoryByDiskSpace,
					Log:      fmt.Sprintf("MountPoint: %v Use: %v% ", mountPoint, disk.SpaceUsage.Percent),
					Status:   0,
				}).Create(global.PanelDB)
				if err != nil {
					return err
				}
				//通知
				err = (&NotifyService{}).SendNotify(
					EventConfig.DiskSpace.NotifyId,
					constant.NotifyLevelWarning,
					fmt.Sprintf("Panel[%v] DiskSpace threshold", global.Config.System.PanelName),
					fmt.Sprintf("MountPoint: %v Use: %v% ", mountPoint, disk.SpaceUsage.Percent))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// DiskInodeEvent DiskInode事件
func (s *MonitorService) DiskInodeEvent(diskList *[]response.Disk) error {
	if !EventConfig.DiskInode.Status || EventConfig.DiskInode.NotifyId == 0 {
		return nil
	}
	for _, mountPoint := range EventConfig.DiskInode.MountPoint {
		for _, disk := range *diskList {
			if disk.MountPoint == mountPoint && EventConfig.DiskInode.Threshold > disk.SpaceUsage.Percent {
				//查询对应事件状态为未处理的数量
				count, err := (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"category": constant.MonitorEventCategoryByDiskInode, "status": 0, "FIXED": fmt.Sprintf("create_time > %d", util.GetDayTimestampZ())})
				if err != nil {
					return err
				}
				if EventConfig.DiskInode.NotifyCount <= count {
					return nil
				}
				//插入事件
				_, err = (&model.MonitorEvent{
					Category: constant.MonitorEventCategoryByDiskInode,
					Log:      fmt.Sprintf("MountPoint: %v Use: %v% ", mountPoint, disk.SpaceUsage.Percent),
					Status:   0,
				}).Create(global.PanelDB)
				if err != nil {
					return err
				}
				//通知
				err = (&NotifyService{}).SendNotify(
					EventConfig.DiskInode.NotifyId,
					constant.NotifyLevelWarning,
					fmt.Sprintf("Panel[%v] DiskInode threshold", global.Config.System.PanelName),
					fmt.Sprintf("MountPoint: %v Use: %v% ", mountPoint, disk.SpaceUsage.Percent))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// LoginPanelEvent LoginPanel事件
func (s *MonitorService) LoginPanelEvent(ip string) error {
	//如果通知状态为开启并且设置了通知id
	if EventConfig.LoginPanel.Status && EventConfig.LoginPanel.NotifyId != 0 {
		//查询对应事件状态为未处理的数量
		count, err := (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"category": constant.MonitorEventCategoryByLoginPanel, "status": 0, "FIXED": fmt.Sprintf("create_time > %d", util.GetDayTimestampZ())})
		if err != nil {
			return err
		}
		if EventConfig.LoginPanel.NotifyCount <= count {
			return nil
		}
		//插入事件
		_, err = (&model.MonitorEvent{
			Category: constant.MonitorEventCategoryByLoginPanel,
			Log:      ip,
			Status:   0,
		}).Create(global.PanelDB)
		if err != nil {
			return err
		}
		//通知
		err = (&NotifyService{}).SendNotify(
			EventConfig.LoginPanel.NotifyId,
			constant.NotifyLevelWarning,
			fmt.Sprintf("Panel[%v] Login Panel", global.Config.System.PanelName),
			fmt.Sprintf("Login IP: %s ", ip))
		if err != nil {
			return err
		}
	}
	return nil
}

// SslExpirationTimeEvent SSLExpirationTime事件
func (s *MonitorService) SslExpirationTimeEvent(projectName, domain string, day int) error {
	//如果通知状态为开启并且设置了通知id
	if EventConfig.SslExpirationTime.Status && EventConfig.SslExpirationTime.NotifyId != 0 && EventConfig.SslExpirationTime.Day >= day {
		//查询对应事件状态为未处理的数量
		count, err := (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"category": constant.MonitorEventCategoryBySslExpirationTime, "status": 0, "FIXED": fmt.Sprintf("create_time > %d", util.GetDayTimestampZ())})
		if err != nil {
			return err
		}
		if EventConfig.SslExpirationTime.NotifyCount <= count {
			return nil
		}
		//插入事件
		_, err = (&model.MonitorEvent{
			Category: constant.MonitorEventCategoryBySslExpirationTime,
			Log:      fmt.Sprintf("ProjectName: %v Domain: %v ,expires in %v days ", projectName, domain, day),
			Status:   0,
		}).Create(global.PanelDB)
		if err != nil {
			return err
		}
		//通知
		err = (&NotifyService{}).SendNotify(
			EventConfig.SslExpirationTime.NotifyId,
			constant.NotifyLevelWarning,
			fmt.Sprintf("Panel[%v] Login Panel", global.Config.System.PanelName),
			fmt.Sprintf("ProjectName: %v Domain: %v ,expires in %v days ", projectName, domain, day))
		if err != nil {
			return err
		}
	}
	return nil
}

// ProjectExpirationTimeEvent 项目到期事件
func (s *MonitorService) ProjectExpirationTimeEvent(projectList []string) error {
	//如果通知状态为开启并且设置了通知id
	if !EventConfig.ProjectExpirationTime.Status && EventConfig.ProjectExpirationTime.NotifyId == 0 {
		return nil
	}
	if len(projectList) > 0 {
		//查询对应事件状态为未处理的数量
		count, err := (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"category": constant.MonitorEventCategoryByProjectExpirationTime, "status": 0, "FIXED": fmt.Sprintf("create_time > %d", util.GetDayTimestampZ())})
		if err != nil {
			return err
		}
		if EventConfig.ProjectExpirationTime.NotifyCount <= count {
			return nil
		}
		content := ""
		for _, v := range projectList {
			content += fmt.Sprintf("ProjectName: %v Expired \r\n", v)
		}
		//插入事件
		_, err = (&model.MonitorEvent{
			Category: constant.MonitorEventCategoryByProjectExpirationTime,
			Log:      content,
			Status:   0,
		}).Create(global.PanelDB)
		if err != nil {
			return err
		}
		//通知
		err = (&NotifyService{}).SendNotify(
			EventConfig.ProjectExpirationTime.NotifyId,
			constant.NotifyLevelWarning,
			fmt.Sprintf("Panel[%v] Login Panel", global.Config.System.PanelName),
			content)
		if err != nil {
			return err
		}
	}
	return nil
}

// ServiceStatusEvent 软件服务状态事件
func (s *MonitorService) ServiceStatusEvent(serviceName string) error {

	downServiceName := ""
	switch serviceName {
	case constant.MonitorEventCategoryByServiceNginx:
		if !EventConfig.ServiceNginx.Status && EventConfig.ServiceNginx.NotifyId == 0 {
			return nil
		}
		if !(&ExtensionNginxService{}).IsRunning() {
			downServiceName = constant.MonitorEventCategoryByServiceNginx
		}
	case constant.MonitorEventCategoryByServiceMysql:
		if !EventConfig.ServiceMysql.Status && EventConfig.ServiceMysql.NotifyId == 0 {
			return nil
		}
		if !(&ExtensionMysqlService{}).IsRunning() {
			downServiceName = constant.MonitorEventCategoryByServiceMysql
		}
	case constant.MonitorEventCategoryByServiceRedis:
		if !EventConfig.ServiceRedis.Status && EventConfig.ServiceRedis.NotifyId == 0 {
			return nil
		}
		if !(&ExtensionRedisService{}).IsRunning() {
			downServiceName = constant.MonitorEventCategoryByServiceRedis
		}
	case constant.MonitorEventCategoryByServiceDocker:
		if !EventConfig.ServiceDocker.Status && EventConfig.ServiceDocker.NotifyId == 0 {
			return nil
		}
		if !(&ExtensionDockerService{}).IsRunning() {
			downServiceName = constant.MonitorEventCategoryByServiceDocker
		}
	default:
		return nil
	}
	//如果通知状态为开启并且设置了通知id
	if downServiceName != "" {
		//查询对应事件状态为未处理的数量
		count, err := (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"category": downServiceName, "status": 0, "FIXED": fmt.Sprintf("create_time > %d", util.GetDayTimestampZ())})
		if err != nil {
			return err
		}
		if EventConfig.SslExpirationTime.NotifyCount <= count {
			return nil
		}
		//插入事件
		_, err = (&model.MonitorEvent{
			Category: downServiceName,
			Log:      fmt.Sprintf("%s is down!", downServiceName),
			Status:   0,
		}).Create(global.PanelDB)
		if err != nil {
			return err
		}
		//通知
		err = (&NotifyService{}).SendNotify(
			EventConfig.SslExpirationTime.NotifyId,
			constant.NotifyLevelWarning,
			fmt.Sprintf("Panel[%v] Service Status", global.Config.System.PanelName),
			fmt.Sprintf("%s is down!", downServiceName))
		if err != nil {
			return err
		}
	}
	return nil
}
