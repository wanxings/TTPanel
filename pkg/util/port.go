package util

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// CheckPort 校验端口是否合法
func CheckPort(port string) bool {
	if StrIsEmpty(port) {
		return false
	}
	//正则校验
	reg := regexp.MustCompile(`^([1-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`)
	return reg.MatchString(port)
}

// CheckProjectPort 校验端口是否合法(项目类)
func CheckProjectPort(port string) bool {
	//不能包含常用端口
	if port == "888" || port == "443" || port == "8080" || port == "21" || port == "25" || port == "22" {
		return false
	}
	return CheckPort(port)
}

// CheckPortRange 校验端口范围是否合法
func CheckPortRange(portRange string) bool {
	//使用-分割端口范围
	portRangeArr := strings.Split(portRange, "-")
	if len(portRangeArr) != 2 {
		return false
	}
	for _, port := range portRangeArr {
		if !CheckPort(port) {
			return false
		}
	}
	return true
}

// CheckTcpPort 检查tcp端口是否被占用,true为被占用
func CheckTcpPort(port int) bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("lsof -i:%d", port))
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	if strings.Contains(string(output), "LISTEN") {
		return true
	}
	return false
}

// CheckPortOccupied 检查端口是否被占用,true为被占用
func CheckPortOccupied(protocol string, port int) bool {
	cmdStr := ""
	if protocol == "tcp" || protocol == "udp" {
		cmdStr = fmt.Sprintf("lsof -i %s:%d", protocol, port)
	} else {
		return false
	}
	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	if strings.Contains(string(output), "LISTEN") {
		return true
	}
	return false
}

// CheckConnection 检查tcp udp 是否能连接
func CheckConnection(protocol string, ip string, port string, timeout time.Duration) (bool, error) {
	conn, err := net.DialTimeout(protocol, fmt.Sprintf("%s:%s", ip, port), timeout)
	if err != nil {
		return false, err
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	return true, nil
}

// CheckUdpPort 检查udp端口是否被占用
func CheckUdpPort(port int) bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("lsof -i:%d", port))
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	if strings.Contains(string(output), "LISTEN") {
		return true
	}
	return false
}

//
