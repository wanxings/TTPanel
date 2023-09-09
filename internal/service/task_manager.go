package service

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"os"
	"strconv"
	"strings"
)

type TaskManagerService struct {
}

// ProcessList 获取进程列表
func (m *TaskManagerService) ProcessList() ([]response.ProcessInfoP, error) {
	pros, err := process.Processes()
	if err != nil {
		global.Log.Errorf("ProcessList.process.Processes err: %v", err)
		return nil, err
	}
	var processList []response.ProcessInfoP
	for _, v := range pros {
		username, _ := v.Username()
		status, _ := v.Status()
		cpuPercent, _ := v.CPUPercent()
		//childrenIDs, _ := v.Children()
		name, _ := v.Name()
		memoryPercent, _ := v.MemoryPercent()
		MemoryUse := (uint64)(0)
		memoryInfo, _ := v.MemoryInfo()
		if memoryInfo != nil {
			MemoryUse = memoryInfo.RSS
		}
		ioCounters, _ := v.IOCounters()
		numFDs, _ := v.NumFDs()         //占用文件描述符
		numThreads, _ := v.NumThreads() //占用线程
		ppid, _ := v.Ppid()

		processList = append(processList, response.ProcessInfoP{
			Pid:           v.Pid,
			PPid:          ppid,
			Name:          name,
			UserName:      username,
			Status:        status,
			CpuPercent:    cpuPercent,
			MemoryPercent: memoryPercent,
			MemoryUse:     MemoryUse,
			IOCounters:    ioCounters,
			NumThreads:    numThreads,
			NumFDs:        numFDs,
		})
	}
	return processList, nil
}

// KillProcess 结束进程
func (m *TaskManagerService) KillProcess(pid int32, force bool) (bool, error) {
	if pid == 0 || pid == 1 {
		return false, errors.New("no zuo no die")
	}

	//判断结束信号
	killCmd := "kill "
	if force {
		killCmd += "-9 "
	} else {
		killCmd += "-15 "
	}
	killCmd += fmt.Sprintf("%d ", pid)

	_, err := util.ExecShell(killCmd)
	if err != nil {
		return false, err
	}
	return true, nil
}

// StartupList 获取启动项列表
func (m *TaskManagerService) StartupList() ([]*util.DirInfo, error) {
	var FileList []*util.DirInfo
	filePathList := []string{
		"/etc/rc.local", "/etc/profile", "/etc/inittab", "/etc/rc.sysinit",
	}
	for _, v := range filePathList {
		if !util.PathExists(v) {
			continue
		}

		content, err := util.ReadFileStringBody(v)
		if err != nil {
			continue
		}

		formatC := util.FormatFileContent(content)
		if len(formatC) == 0 {
			continue
		}

		fileMode := util.GetFilePerm(v)
		if util.FileModePermToString(fileMode.Perm()) == "644" {
			continue
		}
		fileDetails := util.GetFileDetails(v)
		fileDetails.DirName = v
		FileList = append(FileList, fileDetails)
	}
	dirPathList := []string{
		"/etc/init.d", "/etc/rc.d",
	}
	for _, v := range dirPathList {
		fileInfoList, _ := os.ReadDir(v)
		for _, c := range fileInfoList {
			filePath := v + "/" + c.Name()
			if !util.PathExists(filePath) {
				continue
			}
			if util.IsDir(filePath) {
				continue
			}

			fileMode := util.GetFilePerm(filePath)
			if util.FileModePermToString(fileMode.Perm()) == "644" {
				continue
			}
			fileDetails := util.GetFileDetails(v)
			fileDetails.DirName = filePath
			FileList = append(FileList, fileDetails)
		}

	}
	return FileList, nil
}

// ServiceList 获取服务列表
func (m *TaskManagerService) ServiceList() ([]*response.SystemServiceInfo, error) {
	initPath := "/etc/init.d/"
	var systemServiceList []*response.SystemServiceInfo
	fileList, _ := os.ReadDir(initPath)
	for _, v := range fileList {
		info, _ := v.Info()
		if util.FileModePermToString(info.Mode().Perm()) == "644" {
			continue
		}
		LeaveList := m.GetRunLeave(v.Name())
		systemServiceList = append(systemServiceList, &response.SystemServiceInfo{
			Name: v.Name(),
			R0:   LeaveList[0],
			R1:   LeaveList[1],
			R2:   LeaveList[2],
			R3:   LeaveList[3],
			R4:   LeaveList[4],
			R5:   LeaveList[5],
			R6:   LeaveList[6],
			Ps:   "未知",
		})

	}
	return systemServiceList, nil
}

// DeleteService 删除服务
func (m *TaskManagerService) DeleteService(serviceName string) (bool, error) {
	initPath := "/etc/init.d/"
	systemPath := "/usr/lib/systemd/system/"
	if util.PathExists(systemPath + serviceName + ".service") {
		global.Log.Errorf("DeleteService->util.PathExists:%s", systemPath+serviceName+".service")
		return false, errors.New("not found")
	}

	_, err := util.ExecShell("service " + serviceName + " stop")
	if err != nil {
		return false, err
	}

	if util.PathExists("/usr/sbin/update-rc.d") {
		_, _ = util.ExecShell("update-rc.d -f " + serviceName + " remove")
	} else if util.PathExists("/sbin/chkconfig") {
		_, _ = util.ExecShell("chkconfig --del " + serviceName)
	} else {
		_, _ = util.ExecShell("rm -f /etc/rc0.d/*" + serviceName)
		_, _ = util.ExecShell("rm -f /etc/rc1.d/*" + serviceName)
		_, _ = util.ExecShell("rm -f /etc/rc2.d/*" + serviceName)
		_, _ = util.ExecShell("rm -f /etc/rc3.d/*" + serviceName)
		_, _ = util.ExecShell("rm -f /etc/rc4.d/*" + serviceName)
		_, _ = util.ExecShell("rm -f /etc/rc5.d/*" + serviceName)
		_, _ = util.ExecShell("rm -f /etc/rc6.d/*" + serviceName)
	}

	if !util.PathExists(initPath + serviceName) {
		err = os.Remove(initPath + serviceName)
		if err != nil {
			global.Log.Errorf("DeleteService.os.Remove err: %v", err)
			return false, err
		}
	}

	return true, nil
}

// SetRunLevel 设置服务开机启动
func (m *TaskManagerService) SetRunLevel(serviceName string, runLevel int) (bool, error) {
	//如果级别等于0或者大于6则返回错误
	if runLevel < 1 || runLevel > 5 {
		return false, errors.New("级别不对，用命令行搞")
	}

	systemPath := "/usr/lib/systemd/system/"
	multiPath := "/etc/systemd/system/multi-user.target.wants/"

	if util.PathExists(systemPath + serviceName + ".service") {
		global.Log.Errorf("SetRunLevel->util.PathExists:%s", systemPath+serviceName+".service")
		return false, errors.New("not found")
	}

	initPath := "/etc/init.d/"
	if !util.PathExists(initPath + serviceName) {
		nowLeave, _ := util.ExecShell("runlevel")
		nowLeave = strings.Replace(nowLeave, "N", "", -1)
		if nowLeave != fmt.Sprintf("%d", runLevel) {
			return false, errors.New(helper.Message("task_manager.SystemctlManagedServiceCannotSetStatusOfOtherRunlevel"))
		}
		action := "enable"
		if !util.PathExists(multiPath + serviceName + ".service") {
			action = "disable"
		}
		_, err := util.ExecShell("systemctl " + action + " " + serviceName + ".service")
		if err != nil {
			return false, err
		}
		return true, nil
	}

	rcPath := "/etc/rc" + fmt.Sprintf("%d", runLevel) + ".d/"
	rcFileList, _ := os.ReadDir(rcPath)
	for _, v := range rcFileList {
		if v.Name()[3:] != serviceName {
			continue
		}
		nPath := rcPath + v.Name()
		action := "S"
		if v.Name()[:1] == "S" {
			action = "K"
		}
		dPath := rcPath + action + v.Name()[1:]
		err := os.Rename(nPath, dPath)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (m *TaskManagerService) GetRunLeave(path string) []bool {
	etcPath := "/etc/"
	var runLeaveList []bool
	//循环七次
	for i := 0; i < 7; i++ {
		isLeave := false
		fileInfoList, _ := os.ReadDir(etcPath + "rc" + fmt.Sprintf("%d", i) + ".d")
		for _, c := range fileInfoList {
			if c.Name()[3:] == path && c.Name()[:1] == "S" {
				isLeave = true
			}
		}
		runLeaveList = append(runLeaveList, isLeave)
	}
	return runLeaveList
}

// ConnectionList 获取连接列表
func (m *TaskManagerService) ConnectionList() (*response.ConnectionListP, error) {
	connectionList := &response.ConnectionListP{}
	connections, err := net.Connections("all")
	if err != nil {
		return nil, err
	}
	for _, v := range connections {

		newProcess, _ := process.NewProcess(v.Pid)
		processName, _ := newProcess.Name()
		connection := &response.ConnectionInfo{}
		connection.ProcessName = processName
		connection.Type = v.Type
		connection.Pid = v.Pid
		connection.Family = v.Family
		connection.Fd = v.Fd
		connection.Laddr = v.Laddr
		connection.Raddr = v.Raddr
		connection.Status = v.Status
		connection.Uids = v.Uids

		connectionList.List = append(connectionList.List, connection)
	}
	//获取统计信息
	stat, err := net.IOCounters(true)
	connectionList.Statistics = stat
	return connectionList, nil
}

// LinuxUserList 获取用户列表
func (m *TaskManagerService) LinuxUserList() (*response.LinuxUserListP, error) {
	userList := &response.LinuxUserListP{}
	userList.List = make([]*response.LinuxUserInfo, 0)
	//读取文件内容
	content, err := util.ReadFileStringBody("/etc/passwd")
	if err != nil {
		return userList, nil
	}

	users := util.FormatFileContent(content)
	for _, v := range users {
		userInfo := &response.LinuxUserInfo{}
		//使用:分割
		item := strings.Split(v, ":")
		if len(item) != 7 {
			continue
		}
		uid, _ := strconv.ParseInt(item[2], 10, 0)
		gid, _ := strconv.ParseInt(item[3], 10, 0)
		userInfo.Username = item[0]
		userInfo.UID = int(uid)
		userInfo.GID = int(gid)
		userInfo.Home = item[5]
		userInfo.Shell = item[6]
		userInfo.Group = util.GetGroupNameByGID(int(gid))

		userList.List = append(userList.List, userInfo)
	}
	return userList, nil
}

// DeleteLinuxUser 删除用户
func (m *TaskManagerService) DeleteLinuxUser(username string) error {
	//判断用户是否存在
	if !util.IsUserExist(username) {
		return errors.New("user does not exist")
	}
	//尝试删除用户所有的进程
	_, _ = util.ExecShell("killall -u " + username)
	_, _ = util.ExecShell("pkill -9 -u " + username)
	//删除用户
	delResult, _ := util.ExecShell("userdel -r " + username)
	if strings.Contains(delResult, "process") {
		//如果有进程表示失败
		return errors.New("delete Linux user failed")
	}
	err := util.DeleteUser(username)
	if err != nil {
		return err
	}
	return nil
}
