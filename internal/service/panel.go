package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/model"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/host"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PanelService struct{}
type LatestVersion struct {
	Version           string `json:"version"`
	Description       string `json:"description"`
	UpdateTime        int64  `json:"update_time"`
	PreReleaseVersion string `json:"pre_release_version"`
}

// OperatePanel 操作面板
func (s *PanelService) OperatePanel(action string) error {
	//stop restart
	switch action {
	case "stop":
		//停止面板
		_, err := util.ExecShell("nohup tt stop > /dev/null 2> /dev/null &")
		if err != nil {
			return err
		}
	case "restart":
		//重启面板
		_, err := util.ExecShell("nohup tt restart > /dev/null 2> /dev/null &")
		if err != nil {
			return err
		}
	default:
		return errors.New("action error")
	}
	return nil
}

// OperateServer 操作服务器
func (s *PanelService) OperateServer(action string) error {
	//stop restart
	switch action {
	case "stop":
		//停止服务器
		_, _ = util.ExecShell("nohup init 0 > /dev/null 2> /dev/null &")
	case "restart":
		//重启服务器
		_, _ = util.ExecShell("nohup init 6 > /dev/null 2> /dev/null &")
	default:
		return errors.New("action error")
	}
	return nil
}

// ExtensionList 扩展列表
func (s *PanelService) ExtensionList() ([]*response.ExtensionsInfoResponse, error) {
	var extensionsList []*response.ExtensionsInfoResponse
	start := time.Now()

	//nginx
	nginxInfo, err := GroupApp.ExtensionNginxServiceApp.Info()
	if err != nil {
		global.Log.Errorf("ExtensionList->ServiceGroupApp.ExtensionNginxServiceApp.Info Error:%s\n", err.Error())
	}
	extensionsList = append(extensionsList, nginxInfo)
	global.Log.Debugf("ExtensionList->ServiceGroupApp.ExtensionNginxServiceApp.Info cost time:%v \n", time.Now().Sub(start).Seconds())
	//php
	phpInfo, err := GroupApp.ExtensionPHPServiceApp.Info()
	if err != nil {
		global.Log.Errorf("ExtensionList->ServiceGroupApp.ExtensionPHPServiceApp.Info Error:%s\n", err.Error())
	}
	extensionsList = append(extensionsList, phpInfo...)
	global.Log.Debugf("ExtensionList->ServiceGroupApp.ExtensionPHPServiceApp.Info cost time:%v \n", time.Now().Sub(start).Seconds())
	//mysql
	mysqlInfo, err := GroupApp.ExtensionMysqlServiceApp.Info()
	if err != nil {
		global.Log.Errorf("ExtensionList->ServiceGroupApp.ExtensionMysqlServiceApp.Info Error:%s\n", err.Error())
	}
	extensionsList = append(extensionsList, mysqlInfo)
	global.Log.Debugf("ExtensionList->ServiceGroupApp.ExtensionMysqlServiceApp.Info cost time:%v \n", time.Now().Sub(start).Seconds())
	//docker
	dockerInfo, err := GroupApp.ExtensionDockerServiceApp.Info()
	if err != nil {
		global.Log.Errorf("ExtensionList->ServiceGroupApp.ExtensionDockerServiceApp.Info Error:%s\n", err.Error())
	}
	extensionsList = append(extensionsList, dockerInfo)
	global.Log.Debugf("ExtensionList->ServiceGroupApp.ExtensionDockerServiceApp.Info cost time:%v \n", time.Now().Sub(start).Seconds())
	//phpmyadmin
	phpmyadminInfo, err := GroupApp.ExtensionPhpmyadminServiceApp.Info()
	if err != nil {
		global.Log.Errorf("ExtensionList->ServiceGroupApp.ExtensionPhpmyadminServiceApp.Info Error:%s\n", err.Error())
	}
	extensionsList = append(extensionsList, phpmyadminInfo)
	global.Log.Debugf("ExtensionList->ServiceGroupApp.ExtensionPhpmyadminServiceApp.Info cost time:%v \n", time.Now().Sub(start).Seconds())
	//redis
	redisInfo, err := GroupApp.ExtensionRedisServiceApp.Info()
	if err != nil {
		global.Log.Errorf("ExtensionList->ServiceGroupApp.ExtensionRedisServiceApp.Info Error:%s\n", err.Error())
	}
	extensionsList = append(extensionsList, redisInfo)
	global.Log.Debugf("ExtensionList->ServiceGroupApp.ExtensionRedisServiceApp.Info cost time:%v \n", time.Now().Sub(start).Seconds())
	//nodejs
	nodejsInfo, err := GroupApp.ExtensionNodejsServiceApp.Info()
	if err != nil {
		global.Log.Errorf("ExtensionList->ServiceGroupApp.ExtensionNodejsServiceApp.Info Error:%s\n", err.Error())
	}
	extensionsList = append(extensionsList, nodejsInfo)
	global.Log.Debugf("ExtensionList->ServiceGroupApp.ExtensionNodejsServiceApp.Info cost time:%v \n", time.Now().Sub(start).Seconds())
	return extensionsList, nil
}

// Base 面板基础信息
func (s *PanelService) Base() (map[string]interface{}, error) {
	var baseInfo = make(map[string]interface{})
	hostInfo, _ := host.Info()
	//获取面板版本
	baseInfo["version"] = global.Version
	//获取面板预发布版本
	baseInfo["pre_release_version"] = global.Config.System.PreReleaseVersion
	//获取面板运行模式
	baseInfo["run_mode"] = global.Config.System.RunMode
	//面板IP
	baseInfo["panel_ip"] = global.Config.System.PanelIP
	//面板名称
	baseInfo["panel_name"] = global.Config.System.PanelName
	//获取Linux主机名称
	baseInfo["host_name"] = hostInfo.Hostname
	//发行版本
	baseInfo["release"] = fmt.Sprintf("%v%v", hostInfo.Platform, hostInfo.PlatformVersion)
	//内核版本
	baseInfo["kernel_version"] = hostInfo.KernelVersion
	//启动时间
	baseInfo["boot_time"] = hostInfo.BootTime
	//uptime
	baseInfo["uptime"] = hostInfo.Uptime
	//内核架构
	baseInfo["kernel_arch"] = hostInfo.KernelArch
	//网站数量
	baseInfo["website_count"], _ = (&model.Project{}).Count(global.PanelDB)
	//数据库数量
	baseInfo["database_count"], _ = (&model.Databases{}).Count(global.PanelDB)
	//计划任务数量
	baseInfo["crontab_count"], _ = (&model.CronTask{}).Count(global.PanelDB)
	//未处理的监控事件数量
	baseInfo["monitor_events_count"], _ = (&model.MonitorEvent{}).Count(global.PanelDB, &model.ConditionsT{"status": 0})
	//docker相关信息
	baseInfo["docker"], _ = GroupApp.ExtensionDockerServiceApp.BaseStatistics()
	//auto_check_update
	baseInfo["auto_check_update"] = global.Config.System.AutoCheckUpdate
	//语言
	baseInfo["language"] = global.Config.System.Language
	return baseInfo, nil
}

// CheckUpdate 检查更新
func (s *PanelService) CheckUpdate() (*LatestVersion, error) {
	node := getDownloadNode(global.Config.System.CloudNodes[0], global.Config.System.CloudNodes)
	if util.StrIsEmpty(node) {
		return nil, errors.New("cloud node is unavailable")
	}
	var latestVersion LatestVersion
	resp, err := http.Get(fmt.Sprintf("%s/latestVersion.php?prv=%s&time=%d", node, global.Config.System.PreReleaseVersion, time.Now().Unix()))
	if err != nil {
		fmt.Println(1)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			global.Log.Errorf("CheckUpdate->Body.Close Error:%s", err.Error())
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(2)
		return nil, err
	}

	if err = json.Unmarshal(body, &latestVersion); err != nil {
		fmt.Println(3)
		return nil, err
	}
	fmt.Println(latestVersion)
	if !util.StrIsEmpty(latestVersion.Version) && util.CompareVersion(latestVersion.Version, global.Version) == 1 {
		return &latestVersion, nil
	} else if latestVersion.PreReleaseVersion == "error" {
		return nil, fmt.Errorf("error:%s", latestVersion.Description)
	} else {
		return nil, nil
	}
}

// Update 更新面板
func (s *PanelService) Update(latestVersion string) (string, error) {
	if global.Version == latestVersion {
		return "", errors.New("the current version is the latest version")
	}

	downloadNode := getDownloadNode(global.Config.System.CloudNodes[0], global.Config.System.CloudNodes)
	execStr := fmt.Sprintf(`cd %s && wget -O update.sh %s/install/src/update_panel_%s.sh && bash update.sh %s %s`,
		global.Config.System.PanelPath+"/data/shell", downloadNode, latestVersion, latestVersion, global.Config.System.PreReleaseVersion)

	logPath := fmt.Sprintf("%s/panel/update_panel_%s.log", global.Config.Logger.RootPath, time.Now().Format("20060102150405"))
	err := os.MkdirAll(path.Dir(logPath), 0777)
	if err != nil {
		return "", err
	}

	go func() {
		cmdStr := fmt.Sprintf("/bin/bash -c '%s' > %s 2>&1", execStr, logPath)
		cmd := exec.Command("/bin/bash", "-c", cmdStr)
		err = cmd.Run()
		if err != nil {
			global.Log.Errorf("UpdatePanel->cmd.Run  Error:%s", err)
		}
	}()
	//logFile, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0666)
	//if err != nil {
	//	return "", err
	//}
	//
	//go func() {
	//	defer func(logFile *os.File) {
	//		_ = logFile.Close()
	//	}(logFile)
	//	_, _ = logFile.WriteString(fmt.Sprintf("Update panel, pre:%s current version: %s, latest version: v%s.\n----------------------------------------\n", global.Config.System.PreReleaseVersion, global.Version, latestVersion))
	//	_, _ = logFile.WriteString("Selecting cloud node. \n")
	//	downloadNode := getDownloadNode(global.Config.System.CloudNodes[0], global.Config.System.CloudNodes)
	//	if util.StrIsEmpty(downloadNode) {
	//		_, _ = logFile.WriteString("ERROR:cloud node is unavailable,Try again \n")
	//		return
	//	}
	//	_, _ = logFile.WriteString(fmt.Sprintf("Selected cloud node: %s \n", downloadNode))
	//	//清理临时目录
	//	_, _ = logFile.WriteString("Cleaning temporary directory:/www/tmp. \n")
	//	_ = os.RemoveAll("/www/tmp")
	//	//下载更新包
	//	_, _ = logFile.WriteString(fmt.Sprintf("Downloading update package from %s. \n", downloadNode))
	//	url := fmt.Sprintf("%s/update/%s/TTPanel_%s_%s.tar.gz", downloadNode, latestVersion, runtime.GOARCH, global.Config.System.PreReleaseVersion)
	//	savePath := fmt.Sprintf("/www/tmp/TTPanel_%s.tar.gz", latestVersion)
	//	_, err = util.DownloadFile(savePath, url, true)
	//	if err != nil {
	//		_, _ = logFile.WriteString(fmt.Sprintf("ERROR:download update package failed, %s \n", err.Error()))
	//		return
	//	}
	//	//解压更新包
	//	_, _ = logFile.WriteString("Unpacking update package. \n")
	//	_, err = util.ExecShell(fmt.Sprintf("tar -zxvf %s -C /www/tmp", savePath))
	//	if err != nil {
	//		_, _ = logFile.WriteString(fmt.Sprintf("ERROR:Unpacking update package: %s \n", err.Error()))
	//		return
	//	}
	//	//删除面板模板文件 Todo:暂时不删除
	//	//_, _ = logFile.WriteString("Deleting panel template files. \n")
	//	//_ = os.RemoveAll("/www/panel/Templates")
	//	//覆盖旧文件
	//	_, _ = logFile.WriteString("Overwriting old files. \n")
	//	shell1, err := util.ExecShell("cp -rf /www/tmp/panel /www")
	//	if err != nil {
	//		_, _ = logFile.WriteString(fmt.Sprintf("cp -r -f /www/tmp/panel /www ERROR:%s %s\n", shell1, err.Error()))
	//		return
	//	}
	//
	//	//删除临时目录
	//	_, _ = logFile.WriteString("Cleaning temporary directory:/www/tmp \n")
	//	_ = os.RemoveAll("/www/tmp")
	//
	//	//更新完成
	//	_, _ = logFile.WriteString("Update Successfully! \n")
	//
	//	//3秒后重启面板
	//	_, _ = logFile.WriteString("Restarting panel in 3 seconds. \n")
	//	time.Sleep(3 * time.Second)
	//
	//	err = s.OperatePanel("restart")
	//	if err != nil {
	//		return
	//	}
	//
	//}()
	return logPath, nil
}

func getDownloadNode(defaultNode string, nodes []string) string {
	fastestNode := ""
	var tmpFile1, tmpFile2 *os.File
	var err error
	if tmpFile1, err = os.CreateTemp("", "net_test1.*"); err != nil {
		fmt.Println("Failed to create tmpFile1")
		return defaultNode
	}
	if tmpFile2, err = os.CreateTemp("", "net_test2.*"); err != nil {
		fmt.Println("Failed to create tmpFile2")
		return defaultNode
	}
	defer func(tmpFile1, tmpFile2 *os.File) {
		_ = tmpFile1.Close()
		_ = tmpFile2.Close()
		_ = os.Remove(tmpFile1.Name())
		_ = os.Remove(tmpFile2.Name())
	}(tmpFile1, tmpFile2)

	for _, node := range nodes {
		nodeCheck := make(chan string)
		go func(node string, nodeCheck chan string) {
			startTime := time.Now()
			resp, err := http.Get(node)
			if err != nil {
				nodeCheck <- fmt.Sprintf("0 %s", node)
				return
			}
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)
			if resp.StatusCode != 200 {
				nodeCheck <- fmt.Sprintf("0 %s", node)
				return
			}
			endTime := time.Now()
			duration := endTime.Sub(startTime)
			nodeCheck <- fmt.Sprintf("%d %s", duration.Milliseconds(), node)
		}(node, nodeCheck)

		select {
		case res := <-nodeCheck:
			if !util.StrIsEmpty(res) {
				// 将响应时间大于等于 1500ms 的节点写入 tmpFile1
				fmt.Println("res:", res, " ")
				if t, _ := strconv.Atoi(strings.Split(res, " ")[0]); t < 140 {
					fmt.Println("time:", t, "")
					_, _ = fmt.Fprintf(tmpFile1, "%s\n", res)
				}
				// 将响应时间小于 100ms 且大于等于 1500ms 的节点写入 tmpFile2
				if t, _ := strconv.Atoi(strings.Split(res, " ")[0]); t > 100 && t <= 3000 {
					fmt.Println("time:", t, "")
					_, _ = fmt.Fprintf(tmpFile2, "%s\n", res)
				}
			}
		case <-time.After(5 * time.Second):
			global.Log.Debugf("getDownloadNode->Timeout for :%s\n", node)
		}
	}

	// 筛选出请求最快的节点
	if fileSize, _ := tmpFile1.Stat(); fileSize.Size() > 0 {
		_, _ = tmpFile1.Seek(0, 0)
		scanner := bufio.NewScanner(tmpFile1)
		var res []string
		for scanner.Scan() {
			res = append(res, scanner.Text())
		}
		sort.Slice(res, func(i, j int) bool {
			return strings.Split(res[i], " ")[0] < strings.Split(res[j], " ")[0]
		})
		fastestNode = strings.Split(res[0], " ")[1]
	} else if fileSize, _ := tmpFile2.Stat(); fileSize.Size() > 0 {
		_, _ = tmpFile2.Seek(0, 0)
		scanner := bufio.NewScanner(tmpFile2)
		var res []string
		for scanner.Scan() {
			res = append(res, scanner.Text())
		}
		sort.Slice(res, func(i, j int) bool {
			return strings.Split(res[i], " ")[0] < strings.Split(res[j], " ")[0]
		})
		fastestNode = strings.Split(res[0], " ")[1]
	}

	// 输出最快节点的 URL 地址和响应时间
	global.Log.Debugf("Fastest node:%s\n", fastestNode)

	return fastestNode
}
