package cmd

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/service"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var hash string

// cronTaskCmd represents the cronTask command
var cronTaskCmd = &cobra.Command{
	Use:   "cronTask",
	Short: "cronTask",
	Long:  `cronTask`,
	Run: func(cmd *cobra.Command, args []string) {
		if util.StrIsEmpty(hash) {
			fmt.Println("Not Found TaskHash")
			return
		}
		//查询该计划任务
		cronTaskGet, err := (&model.CronTask{Hash: hash}).Get(global.PanelDB)
		if err != nil {
			fmt.Println(err)
			return
		}
		if cronTaskGet.ID == 0 {
			fmt.Println("Not Found CronTask By hash: ", hash)
			return
		}

		switch cronTaskGet.Category {
		case constant.CronTaskCategoryByBackupProject:
			echoBackupStartInfo()
			backupProject(cronTaskGet)
			echoBackupEndInfo()
		case constant.CronTaskCategoryByBackupDatabase:
			echoBackupStartInfo()
			backupDatabase(cronTaskGet)
			echoBackupEndInfo()
		case constant.CronTaskCategoryByBackupDir:
			echoBackupStartInfo()
			backupDir(cronTaskGet)
			echoBackupEndInfo()
		case constant.CronTaskCategoryByCutLog:
			echoCutLogStartInfo()
			cutLog(cronTaskGet)
			echoCutLogEndInfo()
		case constant.CronTaskCategoryByRequestUrl:
			echoRequestUrlStartInfo()
			requestUrl(cronTaskGet)
			echoRequestUrlEndInfo()
		default:
			fmt.Println("Not Found CronTaskCategory")
		}

		return
	},
}

//func freeMemory(cronTaskGet *model.CronTask) {
//	fmt.Println("freeMemory")
//}

func requestUrl(cronTaskGet *model.CronTask) {
	//解析任务配置
	var config model.RequestUrlConfig
	err := util.JsonStrToStruct(cronTaskGet.Config, &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	if util.StrIsEmpty(config.Url) {
		fmt.Println("Not Found Url")
		return
	}

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second, // 设置超时时间为 5 秒
	}
	resp, err := client.Get(config.Url)
	if err != nil {
		fmt.Println(err)
		requestUrlNotify(1, config.NotifyType, config.NotifyId, cronTaskGet.Name, config.Url, err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != config.StatusCode {
		//状态码不一致
		requestUrlNotify(1, config.NotifyType, config.NotifyId, cronTaskGet.Name, config.Url, fmt.Sprintf("Expected StatusCode:%d ,Owned StatusCode:%d", config.StatusCode, resp.StatusCode))
	}
	if util.StrIsEmpty(config.Keyword) && util.StrIsEmpty(config.Reg) {
		//不需要匹配
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		requestUrlNotify(1, config.NotifyType, config.NotifyId, cronTaskGet.Name, config.Url, err.Error())
		fmt.Println(err)
		return
	}
	bodyStr := string(body)

	if !util.StrIsEmpty(config.Keyword) {
		if strings.Contains(bodyStr, config.Keyword) {
			//关键词匹配成功
			requestUrlNotify(2, config.NotifyType, config.NotifyId, cronTaskGet.Name, config.Url, config.Keyword)
			return
		} else {
			requestUrlNotify(3, config.NotifyType, config.NotifyId, cronTaskGet.Name, config.Url, config.Keyword)
		}
	}
	if !util.StrIsEmpty(config.Reg) {
		re := regexp.MustCompile(config.Reg) // 正则表达式
		if re.MatchString(bodyStr) {
			//正则匹配成功
			requestUrlNotify(2, config.NotifyType, config.NotifyId, cronTaskGet.Name, config.Url, config.Reg)
		} else {
			requestUrlNotify(3, config.NotifyType, config.NotifyId, cronTaskGet.Name, config.Url, config.Reg)
		}
	}
}
func requestUrlNotify(Type int, notifyTypes []int, notifyID int64, cronTaskName string, url string, message string) {
	var notifyLevel string
	var content string
	switch Type {
	case 1:
		//访问出现错误
		notifyLevel = constant.NotifyLevelWarning
		content = fmt.Sprintf("Task: %v \nRequest URL: %v, error: %v \n", cronTaskName, url, message)
	case 2:
		//匹配成功
		notifyLevel = constant.NotifyLevelSuccess
		content = fmt.Sprintf("Task: %v \nRequest URL: %v \n[%s] match success", cronTaskName, url, message)
	case 3:
		//匹配失败
		notifyLevel = constant.NotifyLevelWarning
		content = fmt.Sprintf("Task: %v \nRequest URL: %v \n[%s] match failed\n", cronTaskName, url, message)
	default:
		return
	}
	for _, notifyType := range notifyTypes {
		if notifyType == Type {
			//发送通知

			err := service.GroupApp.NotifyServiceApp.SendNotify(
				notifyID,
				notifyLevel,
				fmt.Sprintf("Panel[%v] request URL", global.Config.System.PanelName),
				content,
			)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	return
}

func cutLog(cronTaskGet *model.CronTask) {
	//解析任务配置
	var config model.CutLogConfig
	err := util.JsonStrToStruct(cronTaskGet.Config, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	//查询项目信息
	qWhere := &model.ConditionsT{"ORDER": "create_time DESC"}
	if config.Id == 0 {
		qWhere = &model.ConditionsT{"FIXED": " 1=1 "}
	} else {
		qWhere = &model.ConditionsT{"ID": config.Id}
	}
	projectList, _, err := (&model.Project{}).List(global.PanelDB, qWhere, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	//执行切割
	for _, projectInfo := range projectList {
		//global.Config.System.WwwLogPath
		fmt.Printf("|-Project[%v] \n", projectInfo.Name)
		logAccessPath := fmt.Sprintf("%s/%s.log", global.Config.System.WwwLogPath, projectInfo.Name)
		logErrorPath := fmt.Sprintf("%s/%s.error.log", global.Config.System.WwwLogPath, projectInfo.Name)
		logAccessTarName := fmt.Sprintf("%s_access_%s.log.tar.gz", projectInfo.Name, time.Now().Format("20060102_150405"))
		logErrorTarName := fmt.Sprintf("%s_error_%s.log.tar.gz", projectInfo.Name, time.Now().Format("20060102_150405"))
		logAccessTmpName := fmt.Sprintf("%s.log", projectInfo.Name)
		logErrorTmpName := fmt.Sprintf("%s.error.log", projectInfo.Name)
		historyDir := fmt.Sprintf("%s/history_backups/%s", global.Config.System.WwwLogPath, projectInfo.Name)

		//创建历史目录
		_ = os.MkdirAll(filepath.Dir(historyDir), os.ModePerm)

		//切割日志
		if err = tarLog(logAccessPath, historyDir, logAccessTmpName, logAccessTarName); err != nil {
			fmt.Println(err)
			return
		}
		if err = tarLog(logErrorPath, historyDir, logErrorTmpName, logErrorTarName); err != nil {
			fmt.Println(err)
			return
		}

		//开始清理超过数量的日志
		cleanMaxNumberLogFile(historyDir, config.Number)
	}
}

func cleanMaxNumberLogFile(historyDir string, maxNumber int) {
	//取出所有日志文件
	dirEntry, err := os.ReadDir(historyDir)
	if err != nil {
		panic(err)
	}

	var accessFileList []string
	var errorFileList []string
	for _, file := range dirEntry {
		if !file.IsDir() {
			if strings.Contains(file.Name(), "_access_") {
				accessFileList = append(accessFileList, filepath.Join(historyDir, file.Name()))
			}
			if strings.Contains(file.Name(), "_error_") {
				errorFileList = append(errorFileList, filepath.Join(historyDir, file.Name()))
			}
		}
	}
	//排序
	sort.Strings(accessFileList)
	sort.Strings(errorFileList)
	//删除过期日志
	for i := 0; i < len(accessFileList)-maxNumber; i++ {
		fmt.Printf("|---Expired log[%v] \n", accessFileList[i])
		_ = os.Remove(accessFileList[i])
		fmt.Printf("|---Clean log[%v] \n", accessFileList[i])
	}
	for i := 0; i < len(errorFileList)-maxNumber; i++ {
		fmt.Printf("|---Expired log[%v] \n", errorFileList[i])
		_ = os.Remove(errorFileList[i])
		fmt.Printf("|---Clean log[%v] \n", errorFileList[i])
	}
}

func tarLog(sourceFilePath string, targetDir string, logTmpName string, logTarName string) error {
	tarPath := fmt.Sprintf("%s/%s", targetDir, logTarName)
	//临时日志路径
	logTmpPath := fmt.Sprintf("%s/%s", targetDir, logTmpName)
	//判断日志文件是否存在
	if util.PathExists(sourceFilePath) {
		//移动日志
		_ = os.Rename(sourceFilePath, logTmpPath)

	} else {
		_, _ = os.Create(logTmpPath)
	}

	//压缩日志
	_, err := (&service.ExplorerService{}).Compress(
		true,
		targetDir,
		[]string{logTmpPath},
		tarPath,
		"tar",
		"",
	)
	if err != nil {
		return err
	}
	//删除临时日志文件
	_ = os.Remove(logTmpPath)
	//
	fmt.Printf("|---Split log to[%v] \n", tarPath)
	return nil
}

func backupProject(cronTaskGet *model.CronTask) {
	//解析任务配置
	var config model.BackupProjectConfig
	err := util.JsonStrToStruct(cronTaskGet.Config, &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	//查询项目信息
	qWhere := &model.ConditionsT{"ORDER": "create_time DESC"}
	if config.Id == 0 {
		qWhere = &model.ConditionsT{"FIXED": " 1=1 "}
	} else {
		qWhere = &model.ConditionsT{"ID": config.Id}
	}
	projectList, _, err := (&model.Project{}).List(global.PanelDB, qWhere, 0, 0)
	if err != nil {
		fmt.Println(err)
		backupNotify(config.NotifyType, config.NotifyId, cronTaskGet.Name, err.Error())
		return
	}

	//执行备份
	for _, projectInfo := range projectList {
		//备份项目
		err = service.GroupApp.BackupServiceApp.BackupProject(config.StorageId, config.KeepLocalFile, projectInfo, cronTaskGet.ID, config.ExclusionRules)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("|-Retain the latest number of backups[%v] \n", config.Number)
		//清理过期备份文件
		err = service.GroupApp.BackupServiceApp.ClearExpiredBackupFile(config.Number, constant.BackupCategoryByProject, projectInfo.ID, cronTaskGet.ID, config.StorageId)
		if err != nil {
			fmt.Println(err)
			return
		}
		//如果保留本地备份，则清理本地过期备份文件
		if config.StorageId != 0 && config.KeepLocalFile == 1 {
			err = service.GroupApp.BackupServiceApp.ClearExpiredBackupFile(config.Number, constant.BackupCategoryByProject, projectInfo.ID, cronTaskGet.ID, 0)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func backupDir(cronTaskGet *model.CronTask) {
	//解析任务配置
	var config model.BackupDirConfig
	err := util.JsonStrToStruct(cronTaskGet.Config, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	//备份文件夹
	err = service.GroupApp.BackupServiceApp.BackupDir(config.StorageId, config.KeepLocalFile, config.Path, cronTaskGet.ID, config.ExclusionRules, "定时任务备份")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("|-Retain the latest number of backups[%v] \n", config.Number)

	//清理过期备份文件
	err = service.GroupApp.BackupServiceApp.ClearExpiredBackupFile(config.Number, constant.BackupCategoryByDir, 0, cronTaskGet.ID, config.StorageId)
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果保留本地备份，则清理本地过期备份文件
	if config.StorageId != 0 && config.KeepLocalFile == 1 {
		err = service.GroupApp.BackupServiceApp.ClearExpiredBackupFile(config.Number, constant.BackupCategoryByDir, 0, cronTaskGet.ID, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func backupDatabase(cronTaskGet *model.CronTask) {
	//解析任务配置
	var config model.BackupDatabaseConfig
	err := util.JsonStrToStruct(cronTaskGet.Config, &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	//查询数据库信息
	qWhere := &model.ConditionsT{"ORDER": "create_time DESC"}
	if config.Id == 0 {
		qWhere = &model.ConditionsT{"FIXED": " 1=1 "}
	} else {
		qWhere = &model.ConditionsT{"ID": config.Id}
	}
	databaseList, _, err := (&model.Databases{}).List(global.PanelDB, qWhere, 0, 0)
	if err != nil {
		fmt.Println(err)
		backupNotify(config.NotifyType, config.NotifyId, cronTaskGet.Name, err.Error())
		return
	}

	//执行备份
	for _, databaseInfo := range databaseList {
		//备份数据库
		err = service.GroupApp.BackupServiceApp.BackupMysqlDatabase(config.StorageId, config.KeepLocalFile, databaseInfo, cronTaskGet.ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("|-Retain the latest number of backups[%v] \n", config.Number)

		//清理过期备份文件
		err = service.GroupApp.BackupServiceApp.ClearExpiredBackupFile(config.Number, constant.BackupCategoryByDatabase, databaseInfo.ID, cronTaskGet.ID, config.StorageId)
		if err != nil {
			fmt.Println(err)
			return
		}
		//如果保留本地备份，则清理本地过期备份文件
		if config.StorageId != 0 && config.KeepLocalFile == 1 {
			err = service.GroupApp.BackupServiceApp.ClearExpiredBackupFile(config.Number, constant.BackupCategoryByDatabase, databaseInfo.ID, cronTaskGet.ID, 0)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func backupNotify(notifyType int, notifyID int64, cronTaskName string, errMessage string) {
	//判断是否通知
	if notifyType == 1 {
		err := (&service.NotifyService{}).SendNotify(
			notifyID,
			constant.NotifyLevelWarning,
			fmt.Sprintf("Panel[%v] Backup", global.Config.System.PanelName),
			fmt.Sprintf("Task: %v \nBackup failed, error message: %v \n", cronTaskName, errMessage),
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

func echoBackupStartInfo() {
	//输出开始备份信息
	fmt.Println(util.GetCmdDelimiter())
	fmt.Printf("★Start backup [%v] \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(util.GetCmdDelimiter())
}
func echoBackupEndInfo() {
	//输出结束备份信息
	fmt.Println(util.GetCmdDelimiter())
	fmt.Printf("★Backup completed [%v] \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(util.GetCmdDelimiter())
}
func echoCutLogEndInfo() {
	fmt.Println(util.GetCmdDelimiter())
	fmt.Printf("★Cut log completed [%v] \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(util.GetCmdDelimiter())
}

func echoCutLogStartInfo() {
	fmt.Println(util.GetCmdDelimiter())
	fmt.Printf("★Start cut log [%v] \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(util.GetCmdDelimiter())
}
func echoRequestUrlEndInfo() {
	fmt.Println(util.GetCmdDelimiter())
	fmt.Printf("★Request URL completed [%v] \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(util.GetCmdDelimiter())
}

func echoRequestUrlStartInfo() {
	fmt.Println(util.GetCmdDelimiter())
	fmt.Printf("★Start request URL [%v] \n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(util.GetCmdDelimiter())
}

func init() {
	toolsCmd.AddCommand(cronTaskCmd)
	cronTaskCmd.Flags().StringVarP(&hash, "hash", "i", "", "任务hash")
	_ = cronTaskCmd.MarkFlagRequired("hash")
}
