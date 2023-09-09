package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	cp "github.com/otiai10/copy"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type ExplorerService struct {
}

var remarkMap = map[string]string{
	"/www/wwwlogs":                            "nginx日志目录",
	"/www/docker_app":                         "docker应用目录",
	"/www/recycle_bin":                        "默认回收站目录",
	"/www/server":                             "扩展程序安装目录",
	"/www/swap":                               "虚拟内存文件,请勿删除!",
	"/www/wwwroot":                            "项目目录",
	"/www/backup":                             "默认备份目录",
	"/www/panel":                              "面板目录",
	"/www/panel/data":                         "面板数据目录",
	"/www/panel/config":                       "面板配置目录",
	"/www/panel/data/cast":                    "主机SSH录像目录",
	"/www/panel/data/extensions":              "扩展程序目录",
	"/www/panel/data/i18n":                    "语言包目录",
	"/www/panel/data/7z":                      "7z压缩程序目录",
	"/www/panel/data/explorer_favorites.json": "文件管理收藏夹数据文件,勿动",
	"/www/server/data":                        "此为Mysql数据库默认数据目录，请勿删除!",
	"/www/server/mysql":                       "Mysql程序目录",
	"/www/server/redis":                       "Redis程序目录",
	"/www/server/nvm":                         "PM2/NVM/NPM程序目录",
	"/www/server/pass":                        "网站BasicAuth认证密码存储目录",
	"/www/server/speed":                       "网站加速数据目录",
	"/www/server/docker":                      "Docker插件程序与数据目录",
	"/www/server/total":                       "网站监控报表数据目录",
	"/www/server/ttwaf":                       "Nginx防火墙数据目录",
	"/www/server/phpmyadmin":                  "phpMyAdmin程序目录",
	"/www/server/stop":                        "网站停用页面目录,请勿删除!",
	"/www/server/nginx":                       "Nginx程序目录",
	"/www/server/cron":                        "计划任务脚本与日志目录",
	"/www/server/php":                         "PHP目录,所有PHP版本的解释器都在此目录下",
	"/proc":                                   "系统进程目录",
	"/dev":                                    "系统设备目录",
	"/sys":                                    "系统调用目录",
	"/tmp":                                    "系统临时文件目录",
	"/var/log":                                "系统日志目录",
	"/var/run":                                "系统运行日志目录",
	"/var/spool":                              "系统队列目录",
	"/var/lock":                               "系统锁定目录",
	"/var/mail":                               "系统邮件目录",
	"/mnt":                                    "系统挂载目录",
	"/media":                                  "系统多媒体目录",
	"/dev/shm":                                "系统共享内存目录",
	"/lib":                                    "系统动态库目录",
	"/lib64":                                  "系统动态库目录",
	"/lib32":                                  "系统动态库目录",
	"/usr/lib":                                "系统动态库目录",
	"/usr/lib64":                              "系统动态库目录",
	"/usr/local/lib":                          "系统动态库目录",
	"/usr/local/lib64":                        "系统动态库目录",
	"/usr/local/libexec":                      "系统动态库目录",
	"/usr/local/sbin":                         "系统脚本目录",
	"/usr/local/bin":                          "系统脚本目录",
}

// GetDir 获取目录列表
func (s *ExplorerService) GetDir(query string, path string, sortKey string, reverse bool, offset int, limit int) (map[string]interface{}, error) {
	fileInfoList, err := os.ReadDir(path)
	if err != nil {
		global.Log.Errorf("GetDir->os.ReadDir  path:%s ,Error:%s", err, err)
		return nil, err
	}
	//var list []map[string]interface{}
	list := make(map[string]interface{})
	var DirList []*util.DirInfo
	var FileList []*util.DirInfo
	l := 0
	o := 0
	for _, d := range fileInfoList {
		//如果是搜索行为并且未匹配则跳过
		if !util.StrIsEmpty(query) && strings.IndexAny(strings.ToLower(d.Name()), strings.ToLower(query)) == -1 {
			continue
		}
		o++
		if l >= limit {
			continue
		}
		if o <= offset {
			continue
		}
		fileDetails := util.GetFileDetails(path + "/" + d.Name())
		//获取文件备注
		fileDetails.Remark = s.GetRemark(path + "/" + d.Name())
		if fileDetails.IsDir {
			DirList = append(DirList, fileDetails)
		} else {
			FileList = append(FileList, fileDetails)
		}
		l++
	}
	if !util.StrIsEmpty(sortKey) {
		list["dirs"] = sortedDir(DirList, sortKey, reverse)
		list["files"] = sortedDir(FileList, sortKey, reverse)
	} else {
		list["dirs"] = DirList
		list["files"] = FileList
	}
	list["total_rows"] = o

	return list, nil
}

// SearchDir 搜索目录列表
func (s *ExplorerService) SearchDir(query string, rootPath string, sortKey string, reverse bool) (map[string]interface{}, error) {
	var DirList []*util.DirInfo
	var FileList []*util.DirInfo
	l := 0
	limit := 2000
	list := make(map[string]interface{})
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.Contains(strings.ToLower(d.Name()), strings.ToLower(query)) {
			return nil
		} else {
			if l >= limit {
				return errors.New("limit")
			}
			fileDetails := util.GetFileDetails(path)
			fileDetails.DirName = strings.TrimPrefix(path, rootPath)
			fileDetails.DirName = strings.TrimPrefix(fileDetails.DirName, "/")
			if fileDetails.IsDir {
				DirList = append(DirList, fileDetails)
			} else {
				FileList = append(FileList, fileDetails)
			}
			l++
		}
		return nil
	})
	if err != nil && fmt.Sprintf("%s", err) != "limit" {
		global.Log.Errorf("SearchDir->filepath.Walk returned %s", err.Error())
		return nil, err
	}
	if !util.StrIsEmpty(sortKey) {
		list["dirs"] = sortedDir(DirList, sortKey, reverse)
		list["files"] = sortedDir(FileList, sortKey, reverse)
	} else {
		list["dirs"] = DirList
		list["files"] = FileList
	}
	list["total_rows"] = l
	return list, nil
}

// ReadFile 读取文件内容
func (s *ExplorerService) ReadFile(path string) (string, error) {
	//检查path是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", errors.New("does not exist")
	}
	if fileInfo.Size() > 1024*1024*5 {
		return "", errors.New("File size exceeds 5MB ")
	}
	return util.ReadFileStringBody(path)
}

// SaveFileBody 保存文件内容
func (s *ExplorerService) SaveFileBody(typeField int, path string, body string) error {
	//Todo:目的是解决shell文件换行问题，可能有问题，留坑
	body = strings.ReplaceAll(body, "\r\n", "\n")
	switch typeField {
	case constant.SaveFileTypeFieldByNormal:
		//尝试保存文件副本
		err := s.SaveFileHistory(path)
		if err != nil {
			return errors.New(fmt.Sprintf("SaveFileHistory Error:%s", err.Error()))
		}
		return s.WriteFile(path, body)
	case constant.SaveFileTypeFieldByNginxConf:
		//备份文件
		backupFilePath := fmt.Sprintf("%s.%d.bak", path, time.Now().Unix())
		err := util.CopyFile(path, backupFilePath)
		if err != nil {
			return err
		}
		err = s.WriteFile(path, body)
		if err != nil {
			return err
		}
		//检查配置
		err = (&ExtensionNginxService{}).CheckConfig()
		if err != nil {
			//出现错误回滚文件
			_, err = util.ExecShell(fmt.Sprintf("rm -rf %s;mv %s %s", path, backupFilePath, path))
			if err != nil {
				return err
			}
			return err
		}
		return nil
	default:
		return errors.New("what are you doing?  ")
	}
}

// WriteFile 写入文件内容
func (s *ExplorerService) WriteFile(path string, content string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		global.Log.Errorf("WriteFile->OpenFile Error:%s\n", err.Error())
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	//_, err = file.Write([]byte(content)) //写入字节切片数据
	_, err = file.WriteString(content) //直接写入字符串数据
	if err != nil {
		global.Log.Errorf("WriteFile->WriteString Error:%s\n", err.Error())
		return err
	}
	return nil
}

// GetAttribute 获取文件（夹）属性
func (s *ExplorerService) GetAttribute(path string) (*response.DirFileAttribute, error) {
	//检查path是否存在
	d, err := os.Lstat(path)
	if err != nil {
		return nil, errors.New("does not exist")
	}
	StatT := d.Sys().(*syscall.Stat_t)

	userId := StatT.Uid

	userInfo, _ := user.LookupId(fmt.Sprint(userId))
	groupId := StatT.Gid
	groupInfo, _ := user.LookupGroupId(fmt.Sprint(groupId))
	specialPermissions, _ := getSpecialPermissions(path)
	fileHistory, err := s.GetFileHistory(path)
	if err != nil {
		return nil, err
	}
	attribute := &response.DirFileAttribute{
		Name:              d.Name(),
		Path:              path,
		Size:              d.Size(),
		IsDir:             d.IsDir(),
		IsLink:            d.Mode()&fs.ModeSymlink != 0,
		Owner:             userInfo.Username,
		OwnerId:           userId,
		Group:             groupInfo.Name,
		GroupId:           groupId,
		SpecialPermission: specialPermissions,
		Perm:              strconv.FormatInt(int64(d.Mode().Perm()), 8),
		PermString:        d.Mode().String(),
		ModTime:           d.ModTime().Unix(),
		StatT:             StatT,
		FileHistoryList:   fileHistory,
	}
	return attribute, nil
}

// BatchChangePermission 批量修改文件（夹）权限和所有者
func (s *ExplorerService) BatchChangePermission(paths []string, userName string, permissions string) error {
	for _, v := range paths {
		//检查path是否存在
		if !util.PathExists(v) {
			return errors.New("does not exist")
		}
		if util.IsDir(v) {
			err := filepath.WalkDir(v, func(path string, d fs.DirEntry, err error) error {
				err = util.ChangePermission(path, userName, permissions)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			err := util.ChangePermission(v, userName, permissions)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// BatchDelete 批量删除文件(夹)
func (s *ExplorerService) BatchDelete(paths []string) error {
	//检查敏感目录
	for _, v := range paths {
		if util.CheckSensitivePath(v) {
			return errors.New(helper.MessageWithMap("explorer.CannotOperateSensitivePath", map[string]any{"Path": v}))
		}
	}

	//检查是否开启回收站
	if global.Config.System.RecycleBin.ExplorerStatus {
		for _, v := range paths {
			//检查是否存在文件夹
			if !util.PathExists(v) {
				return errors.New("does not exist")
			}
			err := GroupApp.RecycleBinServiceApp.MoveToRecycleBin(v)
			if err != nil {
				return err
			}
		}
	} else {
		for _, v := range paths {
			err := os.RemoveAll(v)
			if err != nil {
				global.Log.Errorf("BatchDelete->os.RemoveAll  path:%s ,Error:%s", v, err)
				return err
			}
		}
	}
	return nil
}

// Rename 修改文件(文件夹)名
func (s *ExplorerService) Rename(oldPath string, newPath string) error {
	//检查旧路径是否存在
	if !s.BatchCheckFilesExist([]string{oldPath}) {
		return errors.New("does not exist")
	}
	//检查新路径是否存在
	if s.BatchCheckFilesExist([]string{newPath}) {
		return errors.New("path already exists")
	}
	//检查敏感目录
	if util.CheckSensitivePath(oldPath) {
		return errors.New(helper.MessageWithMap("explorer.CannotOperateSensitivePath", map[string]any{"Path": oldPath}))
	}
	if util.CheckSensitivePath(newPath) {
		return errors.New(helper.MessageWithMap("explorer.CannotOperateSensitivePath", map[string]any{"Path": newPath}))
	}
	err := os.Rename(oldPath, newPath)
	if err != nil {
		global.Log.Errorf("Rename->os.Rename  oldPath:%s newPath:%s ,Error:%s", oldPath, newPath, err)
		return err
	}
	return nil
}

// BatchCheckExistsFiles 批量检查文件(夹)是否存在
func (s *ExplorerService) BatchCheckExistsFiles(initPath string, fileList []string, checkPath string) (*response.BatchCheckExistsFilesP, error) {
	var checkList response.BatchCheckExistsFilesP

	//检查两个路径是否是文件夹
	if !checkIsDir(initPath) {
		return nil, errors.New("path is not a dir")
	}

	for _, fileName := range fileList {
		oldPath := cleanPath(initPath + "/" + fileName)
		newPath := cleanPath(checkPath + "/" + fileName)

		//检查是否重复
		if oldPath == newPath {
			global.Log.Debugf("BatchCheckExistsFiles->repeat-> oldPath:%s,newPath:%s", oldPath, newPath)
			checkList.RepeatFiles = append(checkList.RepeatFiles, &response.ComparedFiles{
				OldFile: util.GetFileDetails(oldPath),
				NewFile: util.GetFileDetails(newPath),
			})
		}

		//检查旧路径是否存在
		if !s.BatchCheckFilesExist([]string{oldPath}) {
			checkList.OldExistFiles = append(checkList.OldExistFiles, oldPath)
		}

		//检查目标路径是否存在同名
		if s.BatchCheckFilesExist([]string{newPath}) {
			checkList.ExistSameNameFiles = append(checkList.ExistSameNameFiles, &response.ComparedFiles{
				OldFile: util.GetFileDetails(oldPath),
				NewFile: util.GetFileDetails(newPath),
			})
		}

		//判断是否是文件夹
		if checkIsDir(oldPath) {
			//检查目标文件夹是否是源文件夹的子文件夹
			if strings.Contains(newPath, oldPath) {
				global.Log.Debugf("BatchCheckExistsFiles->The destination folder is a subfolder of the source folder-> oldPath:%s,newPath:%s", oldPath, newPath)
				checkList.TargetIsSubDirectoryFiles = append(checkList.TargetIsSubDirectoryFiles, &response.ComparedFiles{
					OldFile: util.GetFileDetails(oldPath),
					NewFile: util.GetFileDetails(newPath),
				})
			}
		}

	}

	//检查敏感目录

	return &checkList, nil
}

// BatchCopy 批量复制文件(夹)
func (s *ExplorerService) BatchCopy(param *request.BatchCopyMoveR) error {
	//待修改，可能不需要自定义Options

	//if param.Action == "customize" {
	//	fromReserveFileList := ConvertStrSlice2Map(param.FromReserveFileList)
	//	toReserveFileList := ConvertStrSlice2Map(param.ToReserveFileList)
	//	allReserveFileList := ConvertStrSlice2Map(param.AllReserveFileList)
	//	//for i := 0; i < b.N; i++ {
	//	//	InMap(m, "m")
	//	//}
	//}

	var copyOption cp.Options
	copyOption.Skip = func(srcInfo os.FileInfo, src, dest string) (bool, error) {
		switch param.Action {
		case "replace":
			//覆盖替换则直接复制
			global.Log.Debugf("BatchCopy->copyOption.Skip->replace src: %s dest:%s", src, dest)
			return false, nil
		case "jump":
			//跳过重复文件
			global.Log.Debugf("BatchCopy->copyOption.Skip->jump src: %s dest:%s", src, dest)
			return true, nil
		//case "customize":
		//	//自定义
		//	//path.Base(filePath)
		//
		default:
			//覆盖替换则直接复制
			return false, nil
		}
	}

	copyOption.OnDirExists = func(src, dest string) cp.DirExistsAction {
		switch param.Action {
		case "replace":
			//覆盖替换则直接复制
			global.Log.Debugf("BatchCopy->copyOption.OnDirExists->replace src: %s dest:%s", src, dest)
			return cp.DirExistsAction(1)
		case "jump":
			//跳过重复文件
			global.Log.Debugf("BatchCopy->copyOption.OnDirExists->jump src: %s dest:%s", src, dest)
			return cp.DirExistsAction(0)
		//case "customize":
		//	//自定义
		//	//path.Base(filePath)
		//
		default:
			//覆盖替换则直接复制
			return cp.DirExistsAction(1)
		}
	}
	copyOption.PreserveTimes = true //保留条目的 atime 和 mtime
	copyOption.PreserveOwner = true //保留所有条目的 uid 和 gid
	for _, file := range param.FileList {
		iPath := cleanPath(param.InitPath + "/" + file)
		tPath := cleanPath(param.ToPath + "/" + file)
		switch param.Action {
		case "replace":
			//覆盖替换则直接复制
			_ = cp.Copy(iPath, tPath, copyOption)
		case "jump":
			//跳过重复文件
			continue
		//case "customize":
		//	//自定义
		//	//path.Base(filePath)
		//
		default:
			//覆盖替换则直接复制
			_ = cp.Copy(iPath, tPath, copyOption)
		}

	}
	return nil
}

// BatchMove 批量移动文件（夹）
func (s *ExplorerService) BatchMove(param *request.BatchCopyMoveR) error {
	for _, file := range param.FileList {
		iPath := cleanPath(param.InitPath + "/" + file)
		tPath := cleanPath(param.ToPath + "/" + file)
		switch param.Action {
		case "replace":
			//覆盖替换则直接移动
			_ = os.Rename(iPath, tPath)
		case "jump":
			//跳过重复文件(夹)
			continue
		//case "customize":
		//	//自定义
		//	//path.Base(filePath)
		//
		default:
			//覆盖替换则直接移动
			_ = os.Rename(iPath, tPath)
		}

	}
	//var errs []error
	//for _, path := range paths {
	//	err := os.Rename(path, newPath)
	//	if err != nil {
	//		global.Log.Errorf("BatchMove->os.Rename  path:%s newPath:%s ,Error:%s", path, newPath, err)
	//		errs = append(errs, err)
	//	}
	//}
	//return errs
	return nil
}

// CreateDir 创建文件夹
func (s *ExplorerService) CreateDir(path string) error {
	if err := util.CheckDirName(path); err != nil {
		return err
	}
	if s.BatchCheckFilesExist([]string{path}) {
		return errors.New("path already exists")
	}
	err := os.Mkdir(path, 0755)
	if err != nil {
		global.Log.Errorf("CreateDir->os.Mkdir  path:%s ,Error:%s", path, err.Error())
		return err
	}
	//修改文件夹权限
	_ = util.ChangePermission(path, "www", "755")
	return nil
}

// CreateFile 创建文件
func (s *ExplorerService) CreateFile(path string) error {
	if err := util.CheckDirName(path); err != nil {
		return err
	}
	if s.BatchCheckFilesExist([]string{path}) {
		return errors.New("path already exists")
	}
	_, err := os.Create(path)
	if err != nil {
		global.Log.Errorf("CreateFile->os.Create  path:%s ,Error:%s", path, err.Error())
		return err
	}
	//修改文件夹权限
	_ = util.ChangePermission(path, "www", "755")
	return nil
}

// CreateSymlink 创建符号链接
func (s *ExplorerService) CreateSymlink(oldPath, newPath string) error {
	if err := util.CheckDirName(oldPath); err != nil {
		return err
	}
	if err := util.CheckDirName(newPath); err != nil {
		return err
	}

	err := os.Symlink(oldPath, newPath)
	if err != nil {
		global.Log.Errorf("CreateSymlink->os.Symlink  oldPath:%s newPath:%s,Error:%s", oldPath, newPath, err.Error())
		return err
	}
	//修改文件夹权限
	_ = util.ChangePermission(newPath, "www", "755")
	return nil
}

// CreateDuplicate 创建副本
func (s *ExplorerService) CreateDuplicate(path string) (newPath string, err error) {
	if !util.PathExists(path) {
		return "", errors.New(fmt.Sprintf("path %s not exist", path))
	}
	for i := 1; i < 200; i++ {
		newPath = fmt.Sprintf("%s-副本(%d)%s", strings.TrimSuffix(path, filepath.Ext(path)), i, filepath.Ext(path))
		if util.PathExists(newPath) {
			continue
		} else {
			cmdStr := "cp"
			if util.IsDir(path) {
				cmdStr += " -r"
			}
			cmdStr += fmt.Sprintf(" \"%s\" \"%s\"", path, newPath)
			_, err = util.ExecShell(cmdStr)
			if err != nil {
				return "", err
			}
			break
		}
	}
	return newPath, nil
}

// GetPathSize 获取文件夹大小
func (s *ExplorerService) GetPathSize(path string) (int64, error) {
	var size int64
	//检查是否是文件夹
	if !checkIsDir(path) {
		return 0, errors.New("path is not a dir")
	}
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// sortedDir 排序目录列表
func sortedDir(List []*util.DirInfo, sortKey string, reverse bool) []*util.DirInfo {
	if util.StrIsEmpty(sortKey) {
		return List
	}
	sort.SliceStable(List, func(i, j int) bool {
		switch sortKey {
		case "name":
			if reverse {
				return List[i].DirName > List[j].DirName
			}
			return List[i].DirName < List[j].DirName
		case "size":
			if reverse {
				return List[i].Size > List[j].Size
			}
			return List[i].Size < List[j].Size
		case "time":
			if reverse {
				return List[i].ModTime > List[j].ModTime
			}
			return List[i].ModTime < List[j].ModTime
		case "perm":
			if reverse {
				return List[i].Perm > List[j].Perm
			}
			return List[i].Perm < List[j].Perm
		default:
			return List[i].DirName < List[j].DirName
		}
	})
	return List
}

// CreateLink 创建软连接
func (s *ExplorerService) CreateLink(path string, link string) error {
	err := os.Symlink(path, link)
	if err != nil {
		global.Log.Errorf("CreateLink->os.Symlink  path:%s link:%s,Error:%s", path, link, err)
		return err
	}
	return nil
}

// getSpecialPermissions 获取特殊权限
func getSpecialPermissions(path string) (string, error) {
	lsattrBin := which(constant.LsattrCmd)
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	op, err := exec.Command(lsattrBin, "-d", path).CombinedOutput()
	if err != nil {
		// Cannot get path status, return true so that immutable bit is not reverted
		return "", err
	}
	attrs := strings.Split(string(op), " ")
	if len(attrs) != 2 {
		return "", nil
	}
	//if strings.Contains(attrs[0], "i") {
	//	// fmt.Println("Path %v already set to immutable", path)
	//	return true
	//}

	return attrs[0], nil
}

// SetSpecialPermissions 设置特殊权限
func (s *ExplorerService) SetSpecialPermissions(path string, attr string) error {
	chattrBin := which(constant.ChattrCmd)
	if _, err := os.Stat(path); err == nil {
		cmd := exec.Command(chattrBin, "+"+attr, path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err = cmd.Run(); err != nil {
			return fmt.Errorf("%s +%s failed: %s. Err: %v", chattrBin, attr, stderr.String(), err)
		}
	}
	return nil
}

// RemoveSpecialPermissions 移除特殊权限
func (s *ExplorerService) RemoveSpecialPermissions(path string, attr string) error {
	chattrBin := which(constant.ChattrCmd)
	if _, err := os.Stat(path); err == nil {
		cmd := exec.Command(chattrBin, "-"+attr, path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err = cmd.Run(); err != nil {
			return fmt.Errorf("%s -%s failed: %s. Err: %v", chattrBin, attr, stderr.String(), err)
		}
	}

	return nil
}

// CopyDir 复制文件(夹)
//func copyDir(src string, dst string) error {
//	err := cp.Copy(src, dst)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// GetA7zBinPath 获取7z二进制文件路径
func GetA7zBinPath() string {
	a7zPath := global.Config.System.PanelPath + "/data/7z/7z"
	if util.IsArch64() {
		a7zPath += "_arm64"
	}
	return a7zPath
}

// Compress 压缩文件(夹) 如果isWait为true则等待压缩完成 否则后台运行返回压缩命令执行结果的日志路径 所有路径或者密码不能包含单引号
func (s *ExplorerService) Compress(isWait bool, path string, fileList []string, dst string, Type string, passWord string) (logPath string, err error) {
	var shell string
	logFile := fmt.Sprintf("%s/compress/%d.log", global.Config.Logger.RootPath, time.Now().Unix())
	_ = os.MkdirAll(filepath.Dir(logFile), os.ModePerm)
	//将fileList用空格分隔
	if len(fileList) == 0 {
		return "", errors.New("file list cannot be empty")
	}
	A7zBinPath := GetA7zBinPath()

	if strings.Contains(dst, "'") {
		return "", errors.New(fmt.Sprintf("%s has single quotes", dst))
	}
	dst = fmt.Sprintf("'%s'", dst)
	if strings.Contains(path, "'") {
		return "", errors.New(fmt.Sprintf("%s has single quotes", path))
	}
	path = fmt.Sprintf("'%s'", path)
	if strings.Contains(passWord, "'") {
		return "", errors.New(fmt.Sprintf("%s has single quotes", passWord))
	}
	passWord = fmt.Sprintf("'%s'", passWord)

	var src string
	for _, v := range fileList {
		if strings.Contains(v, "'") {
			return "", errors.New(fmt.Sprintf("%s has single quotes", v))
		}
		src += fmt.Sprintf("'%s' ", v)
	}

	//src := strings.Join(fileList, " ")
	//判断压缩类型
	switch Type {
	case constant.CompressTypeByTarGz:
		shell = fmt.Sprintf("cd %s && tar -zcvf %s %s ", path, dst, src)
	case constant.CompressTypeByZip:
		shell = fmt.Sprintf("cd %s && %s a -tzip %s %s ", path, A7zBinPath, dst, src)
		if !util.StrIsEmpty(passWord) {
			shell += fmt.Sprintf(" -p%s -mhe", passWord)
		}
	case constant.CompressTypeBy7z:
		shell = fmt.Sprintf("cd %s && %s a -t7z %s %s ", path, A7zBinPath, dst, src)
		if !util.StrIsEmpty(passWord) {
			shell += fmt.Sprintf(" -p%s -mhe", passWord)
		}
	case constant.CompressTypeByGz:
		shell = fmt.Sprintf("cd %s && %s a -tgzip %s %s ", path, A7zBinPath, dst, src)
	case constant.CompressTypeByTar:
		shell = fmt.Sprintf("cd %s && tar -cvf %s %s ", path, dst, src)
	case constant.CompressTypeByTarXz:
		shell = fmt.Sprintf("cd %s && tar -Jcvf %s %s ", path, dst, src)
	case constant.CompressTypeByXz:
		shell = fmt.Sprintf("cd %s && %s a -txz %s %s ", path, A7zBinPath, dst, src)
	default:
		return "", errors.New("not supported compression type")
	}
	fmt.Println(shell)
	if !isWait {
		shell += fmt.Sprintf(" > %s 2>&1 &", logFile)
	}
	result, err := util.ExecShellScript(shell)
	if err != nil {
		return "", errors.New(fmt.Sprintf("compress failed, error: %s \n %s", err, result))
	}

	return logFile, nil
}

// Decompress 解压文件(夹) 如果isWait为true则等待压缩完成 否则后台运行返回压缩命令执行结果的日志路径
func (s *ExplorerService) Decompress(isWait bool, filePath string, destPath string, password string) (logPath string, err error) {
	logFile := fmt.Sprintf("%s/compress/%d.log", global.Config.Logger.RootPath, time.Now().Unix())
	_ = os.MkdirAll(filepath.Dir(logFile), os.ModePerm)

	if strings.Contains(filePath, "'") {
		return "", errors.New(fmt.Sprintf("%s has single quotes", filePath))
	}
	filePath = fmt.Sprintf("'%s'", filePath)
	if strings.Contains(destPath, "'") {
		return "", errors.New(fmt.Sprintf("%s has single quotes", destPath))
	}
	destPath = fmt.Sprintf("'%s'", destPath)
	if strings.Contains(password, "'") {
		return "", errors.New(fmt.Sprintf("%s has single quotes", password))
	}
	password = fmt.Sprintf("'%s'", password)

	var cmdStr string
	A7zBinPath := GetA7zBinPath()
	if strings.HasSuffix(filePath, ".tar.gz'") {
		cmdStr = fmt.Sprintf("tar -zxvf %s -C %s ", filePath, destPath)
	} else if strings.HasSuffix(filePath, ".tar'") {
		cmdStr = fmt.Sprintf("tar -xvf %s -C %s ", filePath, destPath)
	} else {
		cmdStr = fmt.Sprintf("%s x %s -o%s ", A7zBinPath, filePath, destPath)
		if !util.StrIsEmpty(password) {
			cmdStr += fmt.Sprintf(" -p%s ", password)
		}
		cmdStr += " -y"
	}
	fmt.Println(cmdStr)
	if !isWait {
		cmdStr += fmt.Sprintf(" > %s 2>&1 &", logFile)
	}
	result, err := util.ExecShellScript(cmdStr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("decompress failed, error: %s \n %s", err, result))
	}
	return logFile, nil
}

// GetLogContent 获取日志内容
func (s *ExplorerService) GetLogContent(logPath string, location string, line uint) (string, error) {
	if !util.IsFile(logPath) {
		return fmt.Sprintf("Log file not found:%s", logPath), nil
	}
	//判断文件类型是否是text
	if !util.IsTextFile(logPath) {
		return "", errors.New("file is not of type text")
	}
	cmdStr := ""
	if location == "head" {
		cmdStr = "head"
	} else if location == "tail" {
		cmdStr = "tail"
	}
	cmdStr += fmt.Sprintf(" -q -n %d %s", line, logPath)
	result, err := util.ExecShell(cmdStr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("get log content failed, error: %s \n %s", err, result))
	}
	return result, nil
}

// ClearLogContent 清空日志内容
func (s *ExplorerService) ClearLogContent(logPath string) error {
	if !util.IsFile(logPath) {
		return errors.New("file not exist")
	}
	_, err := util.ExecShell(fmt.Sprintf("truncate -s 0 %s", logPath))
	if err != nil {
		return err
	}
	return nil
}

// SearchFileContent 搜索文件内容
func (s *ExplorerService) SearchFileContent(param *request.SearchFileContentR) (matchList map[string]map[int]string, err error) {
	//判断目录是否存在
	if !util.PathExists(param.DirPath) {
		return nil, errors.New("dir not exist")
	}
	//判断是否是目录
	if !util.IsDir(param.DirPath) {
		return nil, errors.New("dir is not a directory")
	}
	matchList = make(map[string]map[int]string)
	//判断是否搜索子目录
	if param.ContainsSubdir {
		_ = filepath.WalkDir(param.DirPath, func(path string, file fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			searchResult, err := s.GetSearchFileContentResult(path, file, param)
			if err != nil {
				return err
			}
			if len(searchResult) > 0 {
				global.Log.Debugf("SearchFileContent->path: %s \n", param.DirPath+"/"+file.Name())
				global.Log.Debugf("SearchFileContent->searchResult: %v \n", searchResult)
				matchList[param.DirPath+"/"+file.Name()] = searchResult
			}
			return nil
		})
	} else {
		fileList, err := os.ReadDir(param.DirPath)
		if err != nil {
			return nil, err
		}
		for _, file := range fileList {
			searchResult, err := s.GetSearchFileContentResult(param.DirPath+"/"+file.Name(), file, param)
			if err != nil {
				return nil, err
			}
			if len(searchResult) > 0 {
				matchList[param.DirPath+"/"+file.Name()] = searchResult
			}
		}
	}

	return matchList, nil
}

// GetSearchFileContentResult 获取搜索文件内容结果
func (s *ExplorerService) GetSearchFileContentResult(fullPath string, fileDirEntry os.DirEntry, param *request.SearchFileContentR) (result map[int]string, err error) {
	result = make(map[int]string)
	if fileDirEntry.IsDir() {
		return
	}

	//判断文件后缀
	if len(param.Suffix) > 0 {
		isContinue := true
		for _, v := range param.Suffix {
			if v == "no_suffix" && strings.LastIndex(fileDirEntry.Name(), ".") == -1 {
				isContinue = false
				break
			}
			if strings.HasSuffix(fileDirEntry.Name(), "."+v) {
				isContinue = false
				break
			}
		}
		if isContinue {
			return
		}
	}
	global.Log.Debugf("GetSearchFileContentResult->fullPath: %v \n", fullPath)
	//判断文件大小
	topFileInfo, _ := fileDirEntry.Info()
	if (topFileInfo.Size() > param.MaxSize) || (topFileInfo.Size() < param.MinSize) {
		return
	}
	// 判断文件类型是否为文本类型
	file, err := os.Open(fullPath)
	if err != nil {
		_ = file.Close()
		global.Log.Error(err)
		return
	}
	// 读取文件内容
	scanner := bufio.NewScanner(file)
	line := 0
	for scanner.Scan() {
		line++
		txt := scanner.Text()
		fmt.Println(txt)
		//判断是否包含关键字
		if len(param.KeywordReg) > 0 { //正则匹配
			for _, reg := range param.KeywordReg {
				if !param.CaseSensitive {
					reg = "(?i)" + reg
				}
				re := regexp.MustCompile(reg)
				if re.MatchString(txt) {
					result[line] = txt
				}
			}
		}
		if len(param.KeywordNormal) > 0 { //普通匹配
			for _, keyword := range param.KeywordNormal {
				regexp.QuoteMeta(keyword)
				var pattern string
				if !param.CaseSensitive {
					pattern = "(?i)" + regexp.QuoteMeta(keyword)
				} else {
					pattern = regexp.QuoteMeta(keyword)
				}
				re := regexp.MustCompile(pattern)
				if re.MatchString(txt) {
					result[line] = txt
				}
			}
		}
	}
	if scanner.Err() != nil {
		global.Log.Error(scanner.Err())
		return
	}
	return
}

// RemoteDownload 远程下载文件
func (s *ExplorerService) RemoteDownload(param *request.RemoteDownloadR) (key string, err error) {
	if !param.Replace {
		param.SavePath, err = util.RenameFileIfExists(param.SavePath)
		if err != nil {
			return
		}
	}
	key = string(util.RandStr(10, util.ALL))
	err = helper.AsyncDownloadFile(param.Url, param.SavePath, key)
	if err != nil {
		return
	}
	return
}

// RemoteDownloadProcess 远程下载文件进度
func (s *ExplorerService) RemoteDownloadProcess(key string) (*helper.Process, error) {
	value, ok := global.GoCache.Get(key)
	if !ok {
		global.Log.Errorf("RemoteDownloadProcess -> global.GoCache.Get nil")
		return nil, errors.New("not found this DownloadProcess")
	}
	process := &helper.Process{}
	_ = util.JsonStrToStruct(value.(string), process)
	return process, nil
}

// BatchCopy 批量复制文件(文件夹)
//func (s *ExplorerService) BatchCopy(paths []string, newPath string) []error {
//	if CheckSensitiveDir(paths) {
//		return []error{errors.New("不能复制到敏感目录")}
//	}
//	var errs []error
//	for _, path := range paths {
//		err := copy(path, newPath)
//		if err != nil {
//			global.Log.Errorf("BatchCopy->Copy  path:%s newPath:%s ,Error:%s", path, newPath, err)
//			errs = append(errs, err)
//		}
//	}
//	return errs
//}

// BatchCheckFilesExist 批量检查文件(夹)是否存在
func (s *ExplorerService) BatchCheckFilesExist(paths []string) bool {
	for _, v := range paths {
		_, err := os.Stat(v)
		if err != nil {
			if os.IsNotExist(err) {
				return false
			}
		}
	}
	return true
}

// # 名称输入序列化
// def xssdecode(self,text):
// try:
// cs = {"&quot":'"',"&#x27":"'"}
// for c in cs.keys():
// text = text.replace(c,cs[c])
//
// str_convert = text
// if sys.version_info[0] == 3:
// import html
// text2 = html.unescape(str_convert)
// else:
// text2 = cgi.unescape(str_convert)
// return text2
// except:
// return text
// 获取命令行参数
func which(bin string) string {
	pathList := []string{"/usr/bin", "/sbin", "/usr/sbin", "/usr/local/bin"}
	for _, p := range pathList {
		if _, err := os.Stat(path.Join(p, bin)); err == nil {
			return path.Join(p, bin)
		}
	}
	return bin
}

// 检查是否是目录
func checkIsDir(path string) bool {
	d, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !d.IsDir() {
		return false
	}
	return true
}

// 去除路径中多余的/
func cleanPath(path string) string {
	path = strings.Replace(path, "//", "/", -1)
	return path
}

// BatchSetSpecialPermission 批量设置特殊权限
func (s *ExplorerService) BatchSetSpecialPermission(action string, paths []string, permission string, recursion bool) []error {
	var errs []error
	var cmdStrList []string
	for _, v := range paths {
		//路径是否存在
		if !util.PathExists(v) {
			errs = append(errs, errors.New(fmt.Sprintf("路径%s不存在", v)))
			continue
		}
		//是否敏感目录
		if util.CheckSensitivePath(v) {
			errs = append(errs, errors.New(fmt.Sprintf("路径%s为敏感目录，不能设置或移除特殊权限", v)))
			continue
		}
		//是否是目录
		if util.IsDir(v) && recursion {
			cmdStrList = append(cmdStrList, fmt.Sprintf("chattr %s%s -R %s", action, permission, v))
		} else {
			cmdStrList = append(cmdStrList, fmt.Sprintf("chattr %s%s %s", action, permission, v))
		}
	}
	//执行命令
	for _, cmdStr := range cmdStrList {
		_, err := util.ExecShell(cmdStr)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// SetRemark 设置备注
func (s *ExplorerService) SetRemark(path string, remark string) error {
	if path == "/" {
		return errors.New("根目录不能设置备注,你不可能会遇到这个问题")
	}
	if !util.PathExists(path) {
		return errors.New(fmt.Sprintf("路径%s不存在", path))
	}
	//取路径的md5
	pathMd5 := util.EncodeMD5(path)

	savePath := global.Config.System.PanelPath + "/data/explorer_remark/" + pathMd5

	_ = os.MkdirAll(filepath.Dir(savePath), 0644)
	if remark == "" {
		//删除备注
		_ = os.RemoveAll(savePath)
	} else {
		//创建备注
		err := util.WriteFile(savePath, []byte(remark), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetRemark 获取备注
func (s *ExplorerService) GetRemark(path string) string {
	if path == "/" {
		return ""
	}
	if value, ok := remarkMap[path]; ok {
		return value
	}
	//取路径的md5
	pathMd5 := util.EncodeMD5(path)
	savePath := global.Config.System.PanelPath + "/data/explorer_remark/" + pathMd5

	if util.PathExists(savePath) {
		//读取备注
		remark, err := util.ReadFileStringBody(savePath)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return remark
	} else {
		return ""
	}
}

// FavoritesList 获取收藏列表
func (s *ExplorerService) FavoritesList() (map[string]int64, error) {
	filePath := global.Config.System.PanelPath + "/data/explorer_favorites.json"
	if !util.PathExists(filePath) {
		err := util.WriteFile(filePath, []byte("{}"), 0644)
		if err != nil {
			global.Log.Error("FavoritesList->创建explorer_favorites.json失败", err)
			return nil, err
		}
	}
	//读取收藏json
	favoritesJson, err := util.ReadFileStringBody(filePath)
	if err != nil {
		return nil, err
	}
	var favoritesMap map[string]int64
	err = util.JsonStrToStruct(favoritesJson, &favoritesMap)
	if err != nil {
		return nil, err
	}
	if favoritesMap == nil {
		favoritesMap = make(map[string]int64)
	}
	return favoritesMap, nil
}

// OperateFavorites 操作收藏
func (s *ExplorerService) OperateFavorites(action string, path string) error {
	filePath := global.Config.System.PanelPath + "/data/explorer_favorites.json"
	if !util.PathExists(filePath) {
		err := util.WriteFile(filePath, []byte("{}"), 0644)
		if err != nil {
			global.Log.Error("创建explorer_favorites.json失败", err)
			return err
		}
	}
	//读取收藏json
	favoritesJson, err := util.ReadFileStringBody(filePath)
	if err != nil {
		return err
	}
	var favoritesMap map[string]int64
	err = util.JsonStrToStruct(favoritesJson, &favoritesMap)
	if err != nil {
		return err
	}
	if favoritesMap == nil {
		favoritesMap = make(map[string]int64)
	}
	if action == "add" {
		if _, ok := favoritesMap[path]; ok {
			return errors.New("已经收藏过了")
		}
		favoritesMap[path] = time.Now().Unix()
	} else {
		delete(favoritesMap, path)
	}

	favoritesJson, err = util.StructToJsonStr(favoritesMap)
	if err != nil {
		return err
	}
	err = util.WriteFile(global.Config.System.PanelPath+"/data/explorer_favorites.json", []byte(favoritesJson), 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetFavorites 获取收藏
func (s *ExplorerService) GetFavorites() (map[string]int64, error) {
	//读取收藏json
	favoritesJson, err := util.ReadFileStringBody(global.Config.System.PanelPath + "/data/explorer_favorites.json")
	if err != nil {
		return nil, err
	}
	var favoritesMap map[string]int64
	err = util.JsonStrToStruct(favoritesJson, &favoritesMap)
	if err != nil {
		return nil, err
	}
	return favoritesMap, nil
}

// GenerateDownloadExternalLink 生成外部下载链接
func (s *ExplorerService) GenerateDownloadExternalLink(path string, expireDay int, description string) (*model.ExternalDownload, error) {
	//判断路径是否是文件
	if !util.IsFile(path) {
		return nil, errors.New("not found this file")
	}
	//获取过期时间
	externalDownload := &model.ExternalDownload{
		Token:       string(util.RandStr(15, util.ALL)),
		FilePath:    path,
		Description: description,
		ExpireTime:  util.GetTimestampAfterDay(expireDay),
	}
	return externalDownload.Create(global.PanelDB)
}

// GetExternalDownloadByToken 根据token获取外部下载信息
func (s *ExplorerService) GetExternalDownloadByToken(token string) (*model.ExternalDownload, error) {
	return (&model.ExternalDownload{Token: token}).Get(global.PanelDB)
}

// DownloadExternalLinkList 获取外部下载链接列表
func (s *ExplorerService) DownloadExternalLinkList() ([]*model.ExternalDownload, error) {
	err := (&model.ExternalDownload{}).CleanExpired(global.PanelDB)
	if err != nil {
		return nil, err
	}
	return (&model.ExternalDownload{}).List(global.PanelDB)
}

// DeleteDownloadExternalLink 删除外部下载链接
func (s *ExplorerService) DeleteDownloadExternalLink(id int64) error {
	return (&model.ExternalDownload{ID: id}).Delete(global.PanelDB, &model.ConditionsT{})
}

// GetFileTemporaryDownloadLink 获取文件临时下载链接
func (s *ExplorerService) GetFileTemporaryDownloadLink(path string) (string, error) {
	//判断路径是否是文件
	if !util.IsFile(path) {
		return "", errors.New("not found this file")
	}
	token := string(util.RandStr(15, util.ALL))
	global.GoCache.Set("tmp_d_"+token, path, time.Minute*1)
	return fmt.Sprintf("/ExternalDownload?token=%s&name=%s", token, url.QueryEscape(filepath.Base(path))), nil
}

// SaveFileHistory 接收文件路径，将根据全局设置判断是否保存文件副本
func (s *ExplorerService) SaveFileHistory(path string) error {
	if !global.Config.System.FileHistory.Status {
		//未开启该功能则直接返回
		return nil
	}
	if !util.IsFile(path) {
		//如果是不是文件直接返回
		return nil
	}

	//拼接保存路径，最后加上 -history 避免冲突
	saveDirPath := strings.ReplaceAll(fmt.Sprintf("%s/file_history%s-history", global.Config.System.DefaultBackupDirectory, path), "//", "/")
	saveFileName := fmt.Sprintf("%d", time.Now().Unix())

	//创建文件夹
	_ = os.MkdirAll(saveDirPath, 0755)

	//获取目录下的文件
	filList, err := os.ReadDir(saveDirPath)
	if err != nil {
		return err
	}
	//复制
	cmdStr := fmt.Sprintf("cp \"%s\" \"%s\"", path, filepath.Join(saveDirPath, saveFileName))
	_, err = util.ExecShell(cmdStr)
	if err != nil {
		return err
	}
	difference := len(filList) - global.Config.System.FileHistory.Count
	if difference > 0 { //如果已存在的历史副本大于设定的值，需要删除多出的最旧的文件
		// 按照修改时间排序文件
		sort.Slice(filList, func(i, j int) bool {
			iFileInfo, _ := filList[i].Info()
			jFileInfo, _ := filList[j].Info()
			return iFileInfo.ModTime().Before(jFileInfo.ModTime())
		})

		// 删除最旧的文件
		for i := 0; i < difference && i < len(filList); i++ {
			global.Log.Debugf("SaveFileHistory->RemoveFile:%s", filepath.Join(saveDirPath, filList[i].Name()))
			_ = os.Remove(filepath.Join(saveDirPath, filList[i].Name()))
		}
	}
	return nil
}

// GetFileHistory 接收文件路径，获取保存文件副本
func (s *ExplorerService) GetFileHistory(path string) ([]util.FileHistoryInfo, error) {
	FileHistoryInfoList := make([]util.FileHistoryInfo, 0)
	//拼接保存路径，最后加上 -history 避免冲突
	saveDirPath := strings.ReplaceAll(fmt.Sprintf("%s/file_history%s-history", global.Config.System.DefaultBackupDirectory, path), "//", "/")
	if !util.PathExists(saveDirPath) {
		return nil, nil
	}
	//获取目录下的文件
	fileInfoList, err := os.ReadDir(saveDirPath)
	if err != nil {
		return nil, err
	}
	for _, dirEntry := range fileInfoList {
		//获取文件md5
		md5, err := util.GetFileMD5(filepath.Join(saveDirPath, dirEntry.Name()))
		if err != nil {
			global.Log.Debugf("GetFileHistory->GetFileMD5 ERROR:%s", err.Error())
		}
		fileInfo, _ := dirEntry.Info()
		historyInfo := util.FileHistoryInfo{
			Path:    filepath.Join(saveDirPath, dirEntry.Name()),
			MD5:     md5,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime().Unix(),
		}
		FileHistoryInfoList = append(FileHistoryInfoList, historyInfo)
	}
	return FileHistoryInfoList, nil
}
