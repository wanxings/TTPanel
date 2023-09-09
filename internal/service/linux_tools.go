package service

import (
	"TTPanel/internal/global"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LinuxToolsService struct {
}

// GetDnsConfig 获取dns配置
func (s *LinuxToolsService) GetDnsConfig() ([]string, error) {
	dnsStr, err := util.ReadFileStringBody("/etc/resolv.conf")
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`(nameserver)\s+(.+)`)
	matches := re.FindAllStringSubmatch(dnsStr, -1)

	var dnsInfo []string
	for _, match := range matches {
		fmt.Println(match)
		if match[1] == "nameserver" {
			dnsInfo = append(dnsInfo, match[2])
		}
	}
	return dnsInfo, nil
}

// SaveDnsConfig 保存dns配置
func (s *LinuxToolsService) SaveDnsConfig(dnsList []string) error {
	var dnsStr string
	for _, dns := range dnsList {
		dnsStr += "nameserver " + dns + "\n"
	}
	//如果不存在备份配置文件则进行备份
	if !util.PathExists("/etc/resolv.conf.TTPanel.bak") {
		_, _ = util.ExecShell("cp /etc/resolv.conf /etc/resolv.conf.TTPanel.bak")
	}
	err := util.WriteFile("/etc/resolv.conf", []byte(dnsStr), 0644)
	if err != nil {
		return err
	}
	return nil
}

// TestDnsConfig 测试dns配置
func (s *LinuxToolsService) TestDnsConfig(dnsList []string) error {
	// 创建一个新的`net.Dialer`实例。
	dialer := net.Dialer{
		Timeout: time.Second * 2,
	}

	for _, dnsIP := range dnsList {
		// 尝试连接到DNS服务器。
		if _, err := dialer.Dial("udp", dnsIP+":53"); err != nil {
			return err
		}
		fmt.Println("连接成功")
	}

	return nil
}

// RecoverDnsConfig 恢复dns配置
func (s *LinuxToolsService) RecoverDnsConfig() error {
	//如果存在备份配置文件则进行恢复
	if util.PathExists("/etc/resolv.conf.TTPanel.bak") {
		_, err := util.ExecShell("cp /etc/resolv.conf.TTPanel.bak /etc/resolv.conf")
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("备份文件不存在")
	}
	return nil
}

// GetTimeZoneConfig 获取时区配置
func (s *LinuxToolsService) GetTimeZoneConfig(useMainZone string) (any, error) {
	mainZoneList := []string{"Asia", "Africa", "America", "Antarctica", "Arctic", "Atlantic", "Australia", "Europe", "Indian", "Pacific"}
	var subZoneList []string

	result, err := os.Readlink("/etc/localtime")
	if err != nil {
		return nil, err
	}
	item := strings.Split(result, "/")
	fmt.Println(item)

	subZone := item[len(item)-1]
	mainZone := item[len(item)-2]

	if !util.StrIsEmpty(useMainZone) {
		mainZone = useMainZone
	}

	dirEntryList, err := os.ReadDir("/usr/share/zoneinfo/" + mainZone)
	if err != nil {
		return nil, err
	}
	for _, dirEntry := range dirEntryList {
		subZoneList = append(subZoneList, dirEntry.Name())
	}

	if !util.StrIsEmpty(useMainZone) {
		subZone = subZoneList[0]
	}

	//获取当前时间
	shellResult, err := util.ExecShell("date +\"%Y-%m-%d %H:%M:%S %Z %z\" | tr -d '\\n'")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"zone":           map[string]string{"main_zone": mainZone, "sub_zone": subZone},
		"date":           shellResult,
		"main_zone_list": mainZoneList,
		"sub_zone_list":  subZoneList,
	}, nil
}

// SetTimeZone 设置时区
func (s *LinuxToolsService) SetTimeZone(mainZone, subZone string) error {
	_, err := util.ExecShell("ln -sf /usr/share/zoneinfo/" + mainZone + "/" + subZone + " /etc/localtime")
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	return s.SyncDate()
}

// SyncDate 同步服务器时间
func (s *LinuxToolsService) SyncDate() error {
	var err error
	node := getDownloadNode(global.Config.System.CloudNodes[0], global.Config.System.CloudNodes)
	if util.StrIsEmpty(node) {
		return errors.New("cloud node is unavailable")
	}
	resp, err := http.Get(fmt.Sprintf("%s/time.php", node))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	timeBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(2)
		return err
	}

	newTime, err := strconv.Atoi(string(timeBody))
	if err != nil {
		return err
	}
	newTime -= 28800

	addTime, err := util.ExecShell("date +\"%Y-%m-%d %H:%M:%S %Z %z\"")
	if err != nil {
		return errors.New(fmt.Sprintf("获取时区偏差失败:%s", err.Error()))
	}
	addTime = strings.ReplaceAll(addTime, "\n", "")
	addTimeItem := strings.Split(addTime, " ")
	addTimeStr := addTimeItem[len(addTimeItem)-1]
	add1 := false
	if addTimeStr[0] == '+' {
		add1 = true
	}
	tem1, err := strconv.Atoi(addTimeStr[1 : len(addTimeItem)-1])
	if err != nil {
		return err
	}
	tem2, err := strconv.Atoi(addTimeStr[len(addTimeStr)-2:])
	if err != nil {
		return err
	}
	addV := tem1*3600 + tem2*60
	if add1 {
		newTime += addV
	} else {
		newTime -= addV
	}

	dateStr := time.Unix(int64(newTime), 0).Format("2006-01-02 15:04:05")

	shell, err := util.ExecShell(fmt.Sprintf("date -s \"%s\"", dateStr))
	if err != nil {
		return errors.New(fmt.Sprintf("err:%s,shell:%s", err.Error(), shell))
	}

	return nil
}

// GetNetworkConfig 获取网卡配置 Todo
func (s *LinuxToolsService) GetNetworkConfig() (map[string]interface{}, error) {
	//获取网卡列表
	//s.GetNetCardList()
	return nil, nil
}

// GetNetCardList 获取网卡列表 Todo
func (s *LinuxToolsService) GetNetCardList() ([]string, error) {
	shellResult, err := util.ExecShell("ls /sys/class/net | grep -v lo")
	if err != nil {
		return nil, err
	}
	//换行分割提取网卡列表
	return strings.Split(shellResult, "\n"), nil
}

// GetNetCardIP 获取网卡绑定的IP Todo
func (s *LinuxToolsService) GetNetCardIP(netCard string) ([]string, error) {
	shellResult, err := util.ExecShell("ip addr show " + netCard + " | grep inet | grep -v inet6 | awk '{print $2}' | cut -d '/' -f1")
	if err != nil {
		return nil, err
	}
	return strings.Split(shellResult, "\n"), nil
}

//ip addr show eth0 | grep inet | grep -v inet6 | awk '{print $2}' | cut -d '/' -f1
//ip add |grep eth0| grep inet

// GetHostsConfig 获取host配置
func (s *LinuxToolsService) GetHostsConfig() (map[string]string, error) {
	hostConfig := make(map[string]string, 0)
	hostStr, err := util.ReadFileStringBody("/etc/hosts")
	if err != nil {
		return nil, err
	}
	//以换行分割
	lines := strings.Split(hostStr, "\n")
	for _, line := range lines {
		words := strings.SplitN(line, " ", 2)
		if len(words) == 2 {
			hostConfig[strings.TrimSpace(words[1])] = strings.TrimSpace(words[0])
		}
	}
	return hostConfig, nil
}

// AddHosts 添加host
func (s *LinuxToolsService) AddHosts(domain, ip string) error {
	hostsList, err := s.GetHostsConfig()
	if err != nil {
		return err
	}
	hostsList[domain] = ip
	var hostsStr string
	for k, v := range hostsList {
		hostsStr += v + " " + k + "\n"
	}
	return util.WriteFile("/etc/hosts", []byte(hostsStr), 0644)
}

// RemoveHosts 删除host
func (s *LinuxToolsService) RemoveHosts(domain string) error {
	hostsList, err := s.GetHostsConfig()
	if err != nil {
		return err
	}
	delete(hostsList, domain)
	var hostsStr string
	for k, v := range hostsList {
		hostsStr += v + " " + k + "\n"
	}
	return util.WriteFile("/etc/hosts", []byte(hostsStr), 0644)
}
