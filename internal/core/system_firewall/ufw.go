package system_firewall

import (
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
	"strings"
)

var (
	_ Firewall = (*ufwServant)(nil)
)

type ufwServant struct {
	name        string
	strategyMap map[int]string
	protocolMap map[int]string
}

// PingStatus 检查防火墙是否允许ping，true允许，false禁止
func (s *ufwServant) PingStatus() bool {
	result, err := util.ExecShell("sysctl net.ipv4.icmp_echo_ignore_all")
	if err != nil {
		global.Log.Debugf("PingStatus:sysctl net.ipv4.icmp_echo_ignore_all Error:%s,result:%s", err, result)
		return false
	}
	if strings.Contains(util.ClearStr(result), `=1`) {
		return false
	}
	if strings.Contains(util.ClearStr(result), `=0`) {
		return true
	}
	global.Log.Debugf("PingStatus:result:%s", result)
	return true
}
func (s *ufwServant) Status() bool {
	status, _ := util.ExecShell("ufw status")
	if strings.Contains(status, "inactive") {
		return false
	} else if strings.Contains(status, "active") {
		return true
	} else {
		return false
	}
}
func (s *ufwServant) Name() string {
	return s.name
}

// CreatePortRule 创建端口规则
// port: 端口号，支持单个端口和端口范围，如：80、80:90
// sourceIp: 来源IP，支持单个IP和IP段，如：127.0.0.1 127.0.0.1/10
// protocol: 协议，支持tcp、udp
// strategy: 策略，支持allow、deny
func (s *ufwServant) CreatePortRule(port, sourceIp, protocol string, strategy int) error {
	//校验端口
	if !util.CheckPort(port) && !util.CheckPortRange(port) {
		return errors.New(helper.MessageWithMap("PortIsIllegalOrCommon", map[string]any{"Port": port}))
	}
	//校验来源IP
	if !util.StrIsEmpty(sourceIp) && !util.CheckIP(sourceIp) {
		return errors.New(helper.MessageWithMap("IPFormatError", map[string]any{"IP": sourceIp}))
	}
	//校验策略
	if util.StrIsEmpty(s.strategyMap[strategy]) {
		return errors.New("strategy is empty")
	}
	//校验协议
	if util.StrIsEmpty(protocol) {
		return errors.New("protocol is empty")
	}
	//构造命令
	cmd := fmt.Sprintf("ufw %s ", s.strategyMap[strategy])
	if !util.StrIsEmpty(sourceIp) {
		cmd += fmt.Sprintf("from %s ", sourceIp)
	}
	cmd += fmt.Sprintf("to any port %s proto %s", port, protocol)
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return errors.New(result)
	}
	return nil
}

// Close 关闭防火墙
func (s *ufwServant) Close() error {
	//构造命令
	cmd := "ufw disable"
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return errors.New(result)
	}
	return nil
}

// Open 开启防火墙
func (s *ufwServant) Open() error {
	//构造命令
	cmd := "ufw enable"
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return errors.New(result)
	}
	return nil
}

// DeletePortRule 删除端口规则
// port: 端口号，支持单个端口和端口范围，如：80、80:90
// sourceIp: 来源IP，支持单个IP和IP段，如：
// protocol: 协议，支持tcp、udp
// strategy: 策略，支持allow、deny
func (s *ufwServant) DeletePortRule(port, sourceIp, protocol string, strategy int) error {
	//校验端口
	if !util.CheckPort(port) && !util.CheckPortRange(port) {
		return errors.New(helper.MessageWithMap("PortIsIllegalOrCommon", map[string]any{"Port": port}))
	}
	//校验来源IP
	if !util.StrIsEmpty(sourceIp) && !util.CheckIP(sourceIp) {
		return errors.New(helper.MessageWithMap("IPFormatError", map[string]any{"IP": sourceIp}))
	}
	//校验策略
	if util.StrIsEmpty(s.strategyMap[strategy]) {
		return errors.New("strategy is empty")
	}
	//校验协议
	if util.StrIsEmpty(protocol) {
		return errors.New("protocol is empty")
	}
	//构造命令
	cmd := fmt.Sprintf("ufw delete %s ", s.strategyMap[strategy])
	if !util.StrIsEmpty(sourceIp) {
		cmd += fmt.Sprintf("from %s ", sourceIp)
	}
	cmd += fmt.Sprintf("to any port %s proto %s", port, protocol)
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return errors.New(result)
	}
	return nil
}

// Reload 重载防火墙
func (s *ufwServant) Reload() error {
	//构造命令
	cmd := "ufw reload"
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return errors.New(result)
	}
	return nil
}

// AllowPing 允许ping
func (s *ufwServant) AllowPing() error {
	result, err := util.ExecShell("sysctl net.ipv4.icmp_echo_ignore_all=0")
	if err != nil {
		return errors.New(fmt.Sprintf("AllowPing.sysctl net.ipv4.icmp_echo_ignore_all=0 Error:%s result:%s", err, result))
	}
	return nil
}

// DenyPing 禁止ping
func (s *ufwServant) DenyPing() error {
	result, err := util.ExecShell("sysctl net.ipv4.icmp_echo_ignore_all=1")
	if err != nil {
		return errors.New(fmt.Sprintf("AllowPing.sysctl net.ipv4.icmp_echo_ignore_all=1 Error:%s result:%s", err, result))
	}
	return nil
}

// CreateIPRule 创建IP规则
func (s *ufwServant) CreateIPRule(ip string, strategy int) error {
	//校验来源IP
	if !util.CheckIP(ip) {
		return errors.New(helper.MessageWithMap("IPFormatError", map[string]any{"IP": ip}))
	}
	//校验策略
	if util.StrIsEmpty(s.strategyMap[strategy]) {
		return errors.New("strategy is empty")
	}
	//构造命令
	cmd := fmt.Sprintf("ufw %s from %s", s.strategyMap[strategy], ip)
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return err
	}
	return nil
}

// DeleteIPRule 删除IP规则
func (s *ufwServant) DeleteIPRule(ip string, strategy int) error {
	//校验来源IP
	if !util.CheckIP(ip) {
		return errors.New(helper.MessageWithMap("IPFormatError", map[string]any{"IP": ip}))
	}
	//校验策略
	if util.StrIsEmpty(s.strategyMap[strategy]) {
		return errors.New("strategy is empty")
	}
	//构造命令
	cmd := fmt.Sprintf("ufw delete %s from %s", s.strategyMap[strategy], ip)
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return err
	}
	return nil
}

// CreateForwardRule 创建转发规则
// sourceIp: 来源IP，支持单个IP和IP段，如：
// sourcePort: 来源端口，支持单个端口和端口段，如：
// protocol: 协议，支持tcp、udp、icmp
// targetIp: 目标IP，支持单个IP和IP段，如：
// targetPort: 目标端口，支持单个端口和端口段，如：
func (s *ufwServant) CreateForwardRule(targetIp, protocol string, sourcePort, targetPort int64) error {
	err := CheckForwardRule(targetIp, protocol, sourcePort, targetPort)
	if err != nil {
		return err
	}
	//构造命令
	cmd := `iptables -t nat -A PREROUTING `
	cmd += fmt.Sprintf(`-p %s --dport %d -j DNAT --to %s:%d`, protocol, sourcePort, targetIp, targetPort)
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return err
	}
	return nil
}

// DeleteForwardRule 删除转发规则
// sourceIp: 来源IP，支持单个IP和IP段，如：
// sourcePort: 来源端口，支持单个端口和端口段，如：
// protocol: 协议，支持tcp、udp、icmp
// targetIp: 目标IP，支持单个IP和IP段，如：
// targetPort: 目标端口，支持单个端口和端口段，如：
func (s *ufwServant) DeleteForwardRule(targetIp, protocol string, sourcePort, targetPort int64) error {
	err := CheckForwardRule(targetIp, protocol, sourcePort, targetPort)
	if err != nil {
		return err
	}
	//构造命令
	cmd := `iptables -t nat -D PREROUTING `
	cmd += fmt.Sprintf(`-p %s --dport %d -j DNAT --to %s:%d`, protocol, sourcePort, targetIp, targetPort)
	//执行命令
	if result, err := util.ExecShell(cmd); err != nil {
		global.Log.Errorf("ExecShell-Eroor:%s  Result:%s,SourceError:%s", cmd, result, err.Error())
		return err
	}
	return nil
}
