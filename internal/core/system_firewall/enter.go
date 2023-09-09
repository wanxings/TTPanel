package system_firewall

import (
	"TTPanel/internal/helper"
	"TTPanel/internal/helper/constant"
	"TTPanel/pkg/util"
	"errors"
	"fmt"
)

type Firewall interface {
	Name() string
	Close() error
	Open() error
	Reload() error
	Status() bool
	PingStatus() bool
	AllowPing() error
	DenyPing() error
	CreatePortRule(port, sourceIp, protocol string, strategy int) error
	DeletePortRule(port, sourceIp, protocol string, strategy int) error
	CreateIPRule(ip string, strategy int) error
	DeleteIPRule(ip string, strategy int) error
	CreateForwardRule(targetIp, protocol string, sourcePort, targetPort int64) error
	DeleteForwardRule(targetIp, protocol string, sourcePort, targetPort int64) error

	//Search(user *model.User, q *QueryReq, offset, limit int) (*QueryResp, error)
}

func Ufw() Firewall {
	return &ufwServant{
		name: "Ufw",
		strategyMap: map[int]string{
			constant.SystemFirewallStrategyDeny:  "deny",
			constant.SystemFirewallStrategyAllow: "allow",
		},
		protocolMap: map[int]string{
			constant.SystemFirewallProtocolTCP: "tcp",
			constant.SystemFirewallProtocolUDP: "udp",
		},
	}
}
func Firewalld() Firewall {
	return &firewalldServant{
		name: "Firewalld",
		strategyMap: map[int]string{
			constant.SystemFirewallStrategyDeny:  "drop",
			constant.SystemFirewallStrategyAllow: "accept",
		},
		protocolMap: map[int]string{
			constant.SystemFirewallProtocolTCP: "tcp",
			constant.SystemFirewallProtocolUDP: "udp",
		},
	}
}

func Iptables() Firewall {
	return &iptablesServant{
		name: "Iptables",
		strategyMap: map[int]string{
			constant.SystemFirewallStrategyDeny:  "DROP",
			constant.SystemFirewallStrategyAllow: "ACCEPT",
		},
		protocolMap: map[int]string{
			constant.SystemFirewallProtocolTCP: "tcp",
			constant.SystemFirewallProtocolUDP: "udp",
		},
	}
}

func CheckForwardRule(targetIp, protocol string, sourcePort, targetPort int64) error {
	//校验IP
	if !util.CheckIP(targetIp) {
		return errors.New(helper.MessageWithMap("IPFormatError", map[string]any{"IP": targetIp}))
	}
	//校验端口
	if !util.CheckPort(fmt.Sprintf("%d", sourcePort)) {
		return errors.New(helper.MessageWithMap("PortIsIllegalOrCommon", map[string]any{"Port": sourcePort}))
	}
	if !util.CheckPort(fmt.Sprintf("%d", targetPort)) {
		return errors.New(helper.MessageWithMap("PortIsIllegalOrCommon", map[string]any{"Port": targetPort}))
	}
	if targetIp == "127.0.0.1" && sourcePort == targetPort {
		return errors.New("sourcePort = targetPort ")
	}
	//校验协议
	if !(protocol == "tcp" || protocol == "udp" || protocol == "tcp/udp") {
		return errors.New("only tcp udp tcp/udp")
	}
	return nil
}
