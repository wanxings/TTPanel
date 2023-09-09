package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/pkg/util"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type BackupService struct{}

// Panel 备份面板 Todo: 未完成
func (s *BackupService) Panel(param *request.BackupPanelR) error {
	//err := (&backup.Backup{}).BackupPanelData(param.StorageId, param.KeepLocalFile, 0)
	//if err != nil {
	//	return err
	//}
	return nil
}

// List 备份列表
func (s *BackupService) List(category int, pid int64, offset, limit int) (interface{}, int64, error) {
	where := model.ConditionsT{"ORDER": "create_time DESC"}
	where["category"] = category
	if pid != 0 {
		where["pid"] = pid
	}
	return (&model.Backup{}).List(global.PanelDB, &where, &model.ConditionsT{}, offset, limit)
}

// BackupMysqlDatabase 备份数据库
func (s *BackupService) BackupMysqlDatabase(storageID int64, keepLocalFile int, databaseInfo *model.Databases, CronTaskID int64) (err error) {
	if storageID == 0 {
		keepLocalFile = 1
	}
	//输出信息
	exportMysqlStartTime := time.Now().Unix()
	fmt.Printf("|-database[%s] \n", databaseInfo.Name)

	dbService, err := GroupApp.DatabaseMysqlServiceApp.NewMysqlServiceBySid(databaseInfo.Sid)
	if err != nil {
		return err
	}

	//获取数据库大小
	var size sql.NullInt64
	err = dbService.Raw(fmt.Sprintf("SELECT SUM(data_length + index_length) FROM information_schema.tables WHERE table_schema = '%s'", databaseInfo.Name)).Scan(&size).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New(fmt.Sprintf("|-The specified database[%v] has no data! \n", databaseInfo.Name))
		} else {
			return err
		}
	}
	// 查询数据库字符集
	var charset string
	err = dbService.Raw(fmt.Sprintf("SELECT DEFAULT_CHARACTER_SET_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = '%s'", databaseInfo.Name)).Scan(&charset).Error
	if err != nil {
		return err
	}

	fmt.Printf("|-size[%s] \n", util.ConvertSize(size.Int64))
	fmt.Printf("|-charset[%v] \n", charset)
	_, _, diskFreeSize, _, inodeSize, _ := util.GetDiskSpaceAndInode(global.Config.System.DefaultBackupDirectory)
	fmt.Printf("|-available disk space：[%v]，available Inode is：[%d] \n", util.ConvertSize(diskFreeSize), inodeSize)
	fmt.Printf("|-Start export database[%v] \n", time.Now().Format("2006-01-02 15:04:05"))
	mysqldumpPath := GroupApp.DatabaseMysqlServiceApp.GetMysqldumpPath()
	if util.StrIsEmpty(mysqldumpPath) {
		return errors.New("mysqldump command not found, please check if mysql is installed")
	}
	gzFileName := fmt.Sprintf("%s_cron_%s.sql.gz", databaseInfo.Name, time.Now().Format("20060102_150405"))
	localPath := fmt.Sprintf("%s/database/%s", global.Config.System.DefaultBackupDirectory, gzFileName)
	_ = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	backupCmdStr := ""
	if databaseInfo.Sid == 0 {
		dbPwd, _ := GroupApp.DatabaseMysqlServiceApp.GetRootPwd()
		backupCmdStr += fmt.Sprintf("%s -E -R --default-character-set=%s --force --hex-blob --opt %s -u root -p%s | gzip > %s",
			mysqldumpPath, charset, databaseInfo.Name, dbPwd, localPath,
		)
	} else {
		//查询远程数据库服务器信息
		remoteDBInfo, err := (&model.DatabaseServer{ID: databaseInfo.Sid}).Get(global.PanelDB)
		if err != nil || remoteDBInfo == nil {
			return err
		}
		dbHost := remoteDBInfo.Host
		backupCmdStr += fmt.Sprintf("%s -h %s -P %d -E -R --default-character-set=%s --force --hex-blob --opt %s -u %s -p%s | gzip > %s",
			mysqldumpPath, dbHost, remoteDBInfo.Port, charset, databaseInfo.Name, remoteDBInfo.User, remoteDBInfo.Password, localPath,
		)
	}

	global.Log.Debugf("backupCmdStr:%s", backupCmdStr)
	shell, err := util.ExecShell(backupCmdStr)
	if err != nil {
		return errors.New(fmt.Sprintf("%s \n %s", err.Error(), shell))
	}

	gzFileSize := util.GetFileSize(localPath)
	if gzFileSize < 10 {
		return errors.New(fmt.Sprintf("gzFileSize is too small,size:%d \n ", gzFileSize))
	}
	exportMysqlEndTime := time.Now().Unix()
	fmt.Printf("|-export is complete,use:[%v],size:[%v] \n", util.FormatDuration(exportMysqlEndTime-exportMysqlStartTime), util.ConvertSize(gzFileSize))

	if storageID != 0 {
		//上传到云存储
		cloudPath := fmt.Sprintf("TTPanel/backup/mysql_database/%s", gzFileName)
		err = s.UploadToCloudStorage(storageID, localPath, cloudPath)
		if err != nil {
			return err
		}

		//备份成功，插入云端备份记录
		_, err = (&model.Backup{
			Category:   constant.BackupCategoryByDatabase,
			Pid:        databaseInfo.ID,
			CronTaskID: CronTaskID,
			StorageId:  storageID,
			FileName:   gzFileName,
			FilePath:   cloudPath,
			Size:       gzFileSize,
		}).Create(global.PanelDB)
	}
	//备份成功，插入本地备份记录
	if keepLocalFile == 1 {
		_, err = (&model.Backup{
			Category:   constant.BackupCategoryByDatabase,
			Pid:        databaseInfo.ID,
			CronTaskID: CronTaskID,
			StorageId:  storageID,
			FileName:   gzFileName,
			FilePath:   localPath,
			Size:       gzFileSize,
		}).Create(global.PanelDB)
		fmt.Printf("|-Save LocalPath[%v] \n", localPath)
	} else {
		_ = os.Remove(localPath)
	}
	return nil
}

func (s *BackupService) BackupProject(storageID int64, keepLocalFile int, projectInfo *model.Project, CronTaskID int64, exclusionRules []string) (err error) {
	if storageID == 0 {
		keepLocalFile = 1
	}
	//输出信息
	compressFileStartTime := time.Now().Unix()

	fmt.Printf("|-project[%s] \n", projectInfo.Name)
	fmt.Printf("|-path[%s] \n", projectInfo.Path)
	fmt.Printf("|-size[%s] \n", util.ConvertSize(util.GetFileSize(projectInfo.Path)))
	fmt.Printf("|-exclusion rules[%v] \n", exclusionRules)
	_, _, diskFreeSize, _, inodeSize, _ := util.GetDiskSpaceAndInode(global.Config.System.DefaultBackupDirectory)
	fmt.Printf("|-available disk space：[%v]，available Inode is：[%d] \n", util.ConvertSize(diskFreeSize), inodeSize)
	fmt.Printf("|-Start compress[%s] \n", time.Now().Format("2006-01-02 15:04:05"))
	//备份项目
	tarFileName := fmt.Sprintf("%s_cron_%s.tar.gz", projectInfo.Name, time.Now().Format("20060102_150405"))
	localPath := fmt.Sprintf("%s/project/%s", global.Config.System.DefaultBackupDirectory, tarFileName)
	_ = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	//构造排除命令
	excludeStr := "--exclude='.user.ini'"
	for _, rule := range exclusionRules {
		excludeStr += fmt.Sprintf("--exclude='%s' ", rule)
	}
	cmdStr := fmt.Sprintf("cd %s && tar -czvf %s %s %s", filepath.Dir(projectInfo.Path), localPath, excludeStr, filepath.Base(projectInfo.Path))

	cmdStr += " 2>{err_log} 1> /dev/nul"
	shell, err := util.ExecShell(cmdStr)
	if err != nil {
		return errors.New(fmt.Sprintf("backup project error: %s %s", err.Error(), shell))
	}
	tarFileSize := util.GetFileSize(localPath)
	compressFileEndTime := time.Now().Unix()
	fmt.Printf("|-compress is complete,use:[%v],size:[%v] \n", util.FormatDuration(compressFileEndTime-compressFileStartTime), util.ConvertSize(tarFileSize))

	if storageID != 0 {
		//上传到云存储
		cloudPath := fmt.Sprintf("TTPanel/backup/project/%s", tarFileName)
		err = s.UploadToCloudStorage(storageID, localPath, cloudPath)
		if err != nil {
			return err
		}
		//备份成功，插入云端备份记录
		_, err = (&model.Backup{
			Category:   constant.BackupCategoryByProject,
			Pid:        projectInfo.ID,
			CronTaskID: CronTaskID,
			StorageId:  storageID,
			FileName:   tarFileName,
			FilePath:   cloudPath,
			Size:       tarFileSize,
		}).Create(global.PanelDB)

	}
	//备份成功，插入本地备份记录
	if keepLocalFile == 1 {
		_, err = (&model.Backup{
			Category:   constant.BackupCategoryByProject,
			Pid:        projectInfo.ID,
			CronTaskID: CronTaskID,
			StorageId:  storageID,
			FileName:   tarFileName,
			FilePath:   localPath,
			Size:       tarFileSize,
		}).Create(global.PanelDB)
		fmt.Printf("|-Save LocalPath[%v] \n", localPath)
	} else {
		_ = os.Remove(localPath)
	}

	return nil
}

// BackupDir 备份文件夹
func (s *BackupService) BackupDir(storageID int64, keepLocalFile int, dirPath string, CronTaskID int64, exclusionRules []string, description string) (err error) {
	if !util.IsDir(dirPath) {
		return errors.New("directory does not exist or is not a directory")
	}
	if storageID == 0 {
		keepLocalFile = 1
	}
	//输出信息
	startTime := time.Now().Unix()
	fmt.Printf("|-path[%v] \n", dirPath)
	fmt.Printf("|-size[%v] \n", util.ConvertSize(util.GetFileSize(dirPath)))
	fmt.Printf("|-Exclusion rules[%v] \n", exclusionRules)
	_, _, diskFreeSize, _, inodeSize, _ := util.GetDiskSpaceAndInode(global.Config.System.DefaultBackupDirectory)
	fmt.Printf("|-available disk space：[%v]，available Inode is：[%d] \n", util.ConvertSize(diskFreeSize), inodeSize)
	fmt.Printf("|-Start compress[%s] \n", time.Now().Format("2006-01-02 15:04:05"))
	//备份文件夹
	tarFileName := fmt.Sprintf("%s_cron_%s.tar.gz", filepath.Base(dirPath), time.Now().Format("20060102_150405"))
	localPath := fmt.Sprintf("%s/dir/%s", global.Config.System.DefaultBackupDirectory, tarFileName)
	_ = os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	//构造排除命令
	excludeStr := "--exclude='.user.ini'"
	for _, rule := range exclusionRules {
		excludeStr += fmt.Sprintf("--exclude='%s' ", rule)
	}
	cmdStr := fmt.Sprintf("cd %s && tar -czvf %s %s %s", filepath.Dir(dirPath), localPath, excludeStr, filepath.Base(dirPath))

	cmdStr += " 2>{err_log} 1> /dev/nul"
	shell, err := util.ExecShell(cmdStr)
	if err != nil {
		return errors.New(fmt.Sprintf("backup dir error: %s %s", err.Error(), shell))
	}

	endTime := time.Now().Unix()
	tarFileSize := util.GetFileSize(localPath)
	fmt.Printf("|-compress is complete,use:[%v],size:[%v] \n", util.FormatDuration(endTime-startTime), util.ConvertSize(tarFileSize))

	if storageID != 0 {
		//上传到云存储
		cloudPath := fmt.Sprintf("TTPanel/backup/dir/%s", tarFileName)
		err = s.UploadToCloudStorage(storageID, localPath, cloudPath)
		if err != nil {
			return err
		}

		//备份成功，插入云端备份记录
		_, err = (&model.Backup{
			Category:    constant.BackupCategoryByDir,
			Pid:         0,
			CronTaskID:  CronTaskID,
			Description: description,
			StorageId:   storageID,
			FileName:    tarFileName,
			FilePath:    cloudPath,
			Size:        tarFileSize,
		}).Create(global.PanelDB)
	}
	//备份成功，插入本地备份记录
	if keepLocalFile == 1 {
		//插入本地备份记录
		_, err = (&model.Backup{
			Category:    constant.BackupCategoryByDir,
			Pid:         0,
			CronTaskID:  CronTaskID,
			Description: description,
			StorageId:   0,
			FileName:    tarFileName,
			FilePath:    localPath,
			Size:        tarFileSize,
		}).Create(global.PanelDB)
		fmt.Printf("|-Save LocalPath[%v] \n", localPath)
	} else {
		_ = os.Remove(localPath)
	}

	return nil
}

// UploadToCloudStorage 上传到云存储
func (s *BackupService) UploadToCloudStorage(storageId int64, localPath, cloudPath string) (err error) {
	uploadFileStartTime := time.Now().Unix()
	//查询存储配置
	storageGet, err := (&model.Storage{ID: storageId}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if storageGet.ID == 0 {
		return errors.New(fmt.Sprintf("Not Found Storage By ID: %d \n", storageId))
	}
	var config request.StorageConfigR
	err = util.JsonStrToStruct(storageGet.Config, &config)
	if err != nil {
		return err
	}
	storageCore, err := GroupApp.StorageServiceApp.NewStorageCore(storageGet.Category, &config)
	if err != nil {
		return err
	}
	//
	fmt.Printf("|-Start upload file[%v] \n", time.Now().Format("2006-01-02 15:04:05"))
	tarFileSize := util.GetFileSize(localPath)
	err = storageCore.Upload(localPath, cloudPath)
	if err != nil {
		return err
	}
	uploadFileEndTime := time.Now().Unix()
	fmt.Printf("|-|-File upload is complete, storage: %v, time-consuming: %v, upload file size: %v \n", storageGet.Name, util.FormatDuration(uploadFileEndTime-uploadFileStartTime), util.ConvertSize(tarFileSize))

	return nil
}

// ClearExpiredBackupFile 清理过期备份文件
func (s *BackupService) ClearExpiredBackupFile(maxBackupCount int, backupCategory int, PID int64, cronTaskID int64, storageID int64) error {
	//检查是否超过最大备份数量
	list, backupTotal, err := (&model.Backup{}).List(global.PanelDB, &model.ConditionsT{
		"category":     backupCategory,
		"pid":          PID,
		"cron_task_id": cronTaskID,
		"storage_id":   storageID,
		"ORDER":        "create_time DESC",
	}, &model.ConditionsT{}, 0, 0)
	if err != nil {
		return err
	}
	if int(backupTotal) > maxBackupCount {
		for i := 0; i < int(backupTotal)-maxBackupCount; i++ {
			err = s.DeleteBackupFile(list[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteBackupFile 删除备份文件
func (s *BackupService) DeleteBackupFile(backup *model.Backup) error {
	if backup.StorageId == 0 { //删除本地备份文件
		err := os.Remove(backup.FilePath)
		if err != nil {
			return err
		}
		fmt.Printf("|-Expired backup files have been cleaned up from local disk[%v] \n", backup.FilePath)
	} else { //删除云存储备份文件
		//查询存储配置
		storageGet, err := (&model.Storage{ID: backup.StorageId}).Get(global.PanelDB)
		if err != nil {
			return err
		}
		if storageGet.ID == 0 {
			return errors.New(fmt.Sprintf("Not Found Storage By ID: %d \n", backup.StorageId))
		}
		var config request.StorageConfigR
		err = util.JsonStrToStruct(storageGet.Config, &config)
		if err != nil {
			return err
		}
		storageCore, err := GroupApp.StorageServiceApp.NewStorageCore(storageGet.Category, &config)
		if err != nil {
			return err
		}
		if _, err = storageCore.Delete(backup.FilePath); err != nil {
			return err
		}
		//
		fmt.Printf("|-Expired backup files have been cleaned up from cloud storage[%v][%v] \n", storageGet.Name, backup.FilePath)
	}
	//删除数据库记录
	err := backup.Delete(global.PanelDB, &model.ConditionsT{})
	if err != nil {
		return err
	}
	return nil
}
