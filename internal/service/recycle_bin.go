package service

import (
	"TTPanel/internal/conf"
	"TTPanel/internal/global"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RecycleBinService struct{}

// Config 		获取回收站配置
func (s *RecycleBinService) Config() (conf.RecycleBin, error) {
	return global.Config.System.RecycleBin, nil
}

// SetStatus 		设置回收站状态
func (s *RecycleBinService) SetStatus(explorerStatus bool) error {
	global.Config.System.RecycleBin.ExplorerStatus = explorerStatus
	//开始写入配置文件
	newConfig := global.Config.System
	global.Vp.Set("system", newConfig)
	err := global.Vp.WriteConfig() // 保存配置文件
	if err != nil {
		return err
	}
	return nil
}

// List 		获取回收站列表
func (s *RecycleBinService) List() (map[string]response.RecycleBinInfo, error) {
	var err error
	jsonStr := "{}"
	if util.PathExists(global.Config.System.RecycleBin.Directory) {
		_ = os.MkdirAll(global.Config.System.RecycleBin.Directory, 0600)
	}
	//获取回收站json信息
	jsonPath := fmt.Sprintf("%s/.recycle_bin.json", global.Config.System.RecycleBin.Directory)
	if util.PathExists(jsonPath) {
		jsonStr, err = util.ReadFileStringBody(jsonPath)
		if err != nil {
			return nil, err
		}
	}

	var recycleBinInfoMap map[string]response.RecycleBinInfo
	err = util.JsonStrToStruct(jsonStr, &recycleBinInfoMap)
	if err != nil {
		return nil, err
	}
	return recycleBinInfoMap, nil
}

// MoveToRecycleBin 移动文件到回收站
func (s *RecycleBinService) MoveToRecycleBin(path string) error {
	var err error
	jsonStr := "{}"
	if !util.PathExists(global.Config.System.RecycleBin.Directory) {
		_ = os.MkdirAll(global.Config.System.RecycleBin.Directory, 0600)
	}
	//获取回收站json信息
	jsonPath := fmt.Sprintf("%s/.recycle_bin.json", global.Config.System.RecycleBin.Directory)
	if util.PathExists(jsonPath) {
		jsonStr, err = util.ReadFileStringBody(jsonPath)
		if err != nil {
			return err
		}
	}

	var recycleBinInfoMap map[string]response.RecycleBinInfo
	err = util.JsonStrToStruct(jsonStr, &recycleBinInfoMap)
	if err != nil {
		return err
	}

	//获取路径信息
	pathInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	hash := string(util.RandStr(32, util.ALL))

	recycleBinInfo := response.RecycleBinInfo{
		Name:       filepath.Base(path),
		IsDir:      pathInfo.IsDir(),
		Size:       pathInfo.Size(),
		DeleteTime: time.Now().Unix(),
		SourcePath: path,
	}
	// 移动文件
	newPath := fmt.Sprintf("%s/%s", global.Config.System.RecycleBin.Directory, hash)
	err = os.Rename(path, newPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	recycleBinInfoMap[hash] = recycleBinInfo

	//写入配置
	newJsonStr, err := util.StructToJsonStr(recycleBinInfoMap)
	if err != nil {
		return err
	}
	err = util.WriteFile(jsonPath, []byte(newJsonStr), 0644)
	if err != nil {
		return err
	}
	return nil

}

// RecoveryFile 从回收站恢复文件,返回恢复的文件（夹）名称
func (s *RecycleBinService) RecoveryFile(hash string, cover bool) (string, error) {
	var err error
	jsonStr := "{}"
	//获取回收站json信息
	jsonPath := fmt.Sprintf("%s/.recycle_bin.json", global.Config.System.RecycleBin.Directory)
	if util.PathExists(jsonPath) {
		jsonStr, err = util.ReadFileStringBody(jsonPath)
		if err != nil {
			return "", err
		}
	}

	var recycleBinInfoMap map[string]response.RecycleBinInfo
	err = util.JsonStrToStruct(jsonStr, &recycleBinInfoMap)
	if err != nil {
		return "", err
	}
	if _, ok := recycleBinInfoMap[hash]; !ok {
		return "", fmt.Errorf("hash信息不存在")
	}

	if cover { //覆盖
		//先执行删除操作
		err = os.RemoveAll(recycleBinInfoMap[hash].SourcePath)
		if err != nil {
			return "", err
		}
	}
	// 移动文件
	recycleBinPath := fmt.Sprintf("%s/%s", global.Config.System.RecycleBin.Directory, hash)
	err = os.Rename(recycleBinPath, recycleBinInfoMap[hash].SourcePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//删除json信息
	recoveryName := recycleBinInfoMap[hash].Name
	delete(recycleBinInfoMap, hash)
	//写入配置
	newJsonStr, err := util.StructToJsonStr(recycleBinInfoMap)
	if err != nil {
		return "", err
	}
	err = util.WriteFile(jsonPath, []byte(newJsonStr), 0644)
	if err != nil {
		return "", err
	}
	return recoveryName, nil
}

// DeleteRecoveryFile 从回收站删除文件,返回删除的文件（夹）名称
func (s *RecycleBinService) DeleteRecoveryFile(hash string) (string, error) {
	var err error
	jsonStr := "{}"
	//获取回收站json信息
	jsonPath := fmt.Sprintf("%s/.recycle_bin.json", global.Config.System.RecycleBin.Directory)
	if util.PathExists(jsonPath) {
		jsonStr, err = util.ReadFileStringBody(jsonPath)
		if err != nil {
			return "", err
		}
	}

	var recycleBinInfoMap map[string]response.RecycleBinInfo
	err = util.JsonStrToStruct(jsonStr, &recycleBinInfoMap)
	if err != nil {
		return "", err
	}
	if _, ok := recycleBinInfoMap[hash]; !ok {
		return "", fmt.Errorf("hash信息不存在")
	}
	recycleBinPath := fmt.Sprintf("%s/%s", global.Config.System.RecycleBin.Directory, hash)
	err = os.RemoveAll(recycleBinPath)
	if err != nil {
		return "", err
	}
	deleteName := recycleBinInfoMap[hash].Name
	delete(recycleBinInfoMap, hash)
	//写入配置
	newJsonStr, err := util.StructToJsonStr(recycleBinInfoMap)
	if err != nil {
		return "", err
	}
	err = util.WriteFile(jsonPath, []byte(newJsonStr), 0644)
	if err != nil {
		return "", err
	}
	return deleteName, nil
}

// ClearRecycleBin 清空回收站
func (s *RecycleBinService) ClearRecycleBin() error {
	if !util.PathExists(global.Config.System.RecycleBin.Directory) {
		return fmt.Errorf("回收站目录不存在")
	}
	if !strings.Contains(global.Config.System.RecycleBin.Directory, "/recycle_bin") {
		global.Log.Errorf("回收站目录异常，不允许清空,回收站目录：%s", global.Config.System.RecycleBin.Directory)
		return fmt.Errorf("回收站目录异常，不允许清空")
	}

	cmdStr := fmt.Sprintf("rm -rf %s/*", global.Config.System.RecycleBin.Directory)
	_, err := util.ExecShell(cmdStr)
	if err != nil {
		return err
	}

	err = util.WriteFile(fmt.Sprintf("%s/.recycle_bin.json", global.Config.System.RecycleBin.Directory), []byte("{}"), 0600)
	if err != nil {
		return err
	}
	return nil
}
