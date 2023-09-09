package service

import (
	"TTPanel/internal/core/system_firewall"
	"TTPanel/internal/global"
	"TTPanel/internal/helper"
	"TTPanel/internal/model"
	"TTPanel/internal/model/request"
	"TTPanel/internal/model/response"
	"TTPanel/pkg/util"
	"errors"
	"strconv"
	"strings"
)

type SystemFirewallService struct {
	Core system_firewall.Firewall
}

func (s *SystemFirewallService) New() (*SystemFirewallService, error) {
	var err error
	if _, err = util.ExecShell("which ufw"); err == nil {
		s.Core = system_firewall.Ufw()
	} else if _, err = util.ExecShell("which firewalld"); err == nil {
		s.Core = system_firewall.Firewalld()
	} else if _, err = util.ExecShell("which iptables"); err == nil {
		s.Core = system_firewall.Iptables()
	} else {
		return nil, errors.New(helper.Message("firewalld.NotInstalled"))
	}
	return s, nil
}

// FirewallStatus 获取系统防火墙信息
func (s *SystemFirewallService) FirewallStatus() *response.FirewallInfo {
	//获取ssh状态，
	return &response.FirewallInfo{
		Ping:   s.Core.PingStatus(),
		Status: s.Core.Status(),
		Name:   s.Core.Name(),
	}
}

// BatchCreatePortRule 批量创建端口规则
func (s *SystemFirewallService) BatchCreatePortRule(rules []*request.CreatePortRuleR) error {
	for _, rule := range rules {
		//检查数据库中协议+端口是否已存在
		_, total, err := (&model.FirewallRulePort{}).List(global.PanelDB, &model.ConditionsT{
			"port = ?":     rule.Port,
			"protocol = ?": rule.Protocol,
			"strategy = ?": rule.Strategy,
		}, 0, 0)
		if err != nil {
			return err
		}
		if total > 0 {
			return errors.New(helper.Message("firewalld.SourcePortAlreadyExists"))
		}
		if err := s.Core.CreatePortRule(strconv.Itoa(rule.Port), rule.SourceIp, rule.Protocol, rule.Strategy); err != nil {
			return err
		}
		//重载防火墙
		if err := s.Core.Reload(); err != nil {
			return err
		}
		//保存到数据库
		insertData := &model.FirewallRulePort{
			Port:     strconv.Itoa(rule.Port),
			SourceIp: rule.SourceIp,
			Protocol: rule.Protocol,
			Strategy: rule.Strategy,
			Ps:       rule.Ps,
		}
		_, err = (insertData).Create(global.PanelDB)
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchDeletePortRule 批量删除端口规则
func (s *SystemFirewallService) BatchDeletePortRule(ids []int64) error {
	for _, id := range ids {
		//从数据库中获取规则
		rule, err := (&model.FirewallRulePort{ID: id}).Get(global.PanelDB)
		if err != nil {
			return err
		}
		if rule.ID == 0 {
			return errors.New("not found")
		}
		//删除防火墙规则
		if err = s.Core.DeletePortRule(rule.Port, rule.SourceIp, rule.Protocol, rule.Strategy); err != nil {
			return err
		}
		//重载防火墙
		if err = s.Core.Reload(); err != nil {
			return err
		}
		//从数据库中删除规则
		if err = rule.Delete(global.PanelDB, &model.ConditionsT{}); err != nil {
			return err
		}
	}
	return nil
}

// GetPortRules 获取端口规则列表
func (s *SystemFirewallService) GetPortRules(query string, page, limit int) ([]*response.FirewallRulePortResp, int64, error) {
	//global.Log.Errorf("query: %s, page: %d, limit: %d", query, page, limit)
	//从数据库中获取规则
	rule, total, err := (&model.FirewallRulePort{}).List(global.PanelDB, &model.ConditionsT{"ORDER": "create_time DESC"}, page, limit)
	if err != nil {
		return nil, 0, err
	}
	//转换为返回数据
	var resp []*response.FirewallRulePortResp
	for _, v := range rule {
		resp = append(resp, &response.FirewallRulePortResp{
			Status:           s.GetPortStatus(v.Protocol, v.Port),
			FirewallRulePort: v,
		})
	}
	return resp, total, nil
}

// GetPortStatus 获得端口状态
func (s *SystemFirewallService) GetPortStatus(network, port string) bool {
	if strings.Contains(port, "-") {
		// 如果是端口范围，拆分为两个端口
		ports := strings.Split(port, "-")
		start, _ := strconv.Atoi(ports[0])
		end, _ := strconv.Atoi(ports[1])

		// 检查端口范围
		for i := start; i <= end; i++ {
			// 检查端口是否被占用
			if !util.CheckPortOccupied(network, i) {
				continue
			} else {
				return true
			}
		}
	} else {
		// 如果是单个端口
		p, _ := strconv.Atoi(port)
		// 检查端口是否被占用
		return util.CheckPortOccupied(network, p)
	}
	return false
}

// UpdatePortRule 更新端口规则
func (s *SystemFirewallService) UpdatePortRule(id int64, rule *request.UpdatePortRuleR) error {
	//从数据库中获取规则
	r, err := (&model.FirewallRulePort{ID: id}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if r.ID == 0 {
		return errors.New("not found")
	}
	//删除防火墙规则
	if err = s.Core.DeletePortRule(r.Port, r.SourceIp, r.Protocol, r.Strategy); err != nil {
		return err
	}
	//重载防火墙
	if err = s.Core.Reload(); err != nil {
		return err
	}
	//更新数据库
	r.Port = rule.Port
	r.SourceIp = rule.SourceIp
	r.Protocol = rule.Protocol
	r.Strategy = rule.Strategy
	r.Ps = rule.Ps
	if err = r.Update(global.PanelDB); err != nil {
		return err
	}
	//添加防火墙规则
	if err = s.Core.CreatePortRule(r.Port, r.SourceIp, r.Protocol, r.Strategy); err != nil {
		return err
	}
	//重载防火墙
	if err = s.Core.Reload(); err != nil {
		return err
	}
	return nil
}

// Close 关闭防火墙
func (s *SystemFirewallService) Close() error {
	return s.Core.Close()
}

// Open 开启防火墙
func (s *SystemFirewallService) Open() error {
	return s.Core.Open()
}

// AllowPing 允许ping
func (s *SystemFirewallService) AllowPing() error {
	return s.Core.AllowPing()
}

// DenyPing 禁止ping
func (s *SystemFirewallService) DenyPing() error {
	return s.Core.DenyPing()
}

// BatchCreateIPRule 批量创建IP规则
func (s *SystemFirewallService) BatchCreateIPRule(rules []*request.CreateIPRuleR) error {
	for _, rule := range rules {
		//检查数据库中源IP是否已存在
		_, total, err := (&model.FirewallRuleIp{}).List(global.PanelDB, &model.ConditionsT{
			"ip = ?": rule.SourceIp,
		}, 0, 0)
		if err != nil {
			return err
		}
		if total > 0 {
			return errors.New(helper.Message("firewalld.IPRuleAlreadyExists"))
		}
		if err := s.Core.CreateIPRule(rule.SourceIp, rule.Strategy); err != nil {
			return err
		}
		//重载防火墙
		if err := s.Core.Reload(); err != nil {
			return err
		}
		//保存到数据库
		insertData := &model.FirewallRuleIp{
			Ip:       rule.SourceIp,
			Strategy: rule.Strategy,
			Ps:       rule.Ps,
		}
		_, err = (insertData).Create(global.PanelDB)
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchDeleteIPRule 批量删除IP规则
func (s *SystemFirewallService) BatchDeleteIPRule(ids []int64) error {
	for _, id := range ids {
		//从数据库中获取规则
		r, err := (&model.FirewallRuleIp{ID: id}).Get(global.PanelDB)
		if err != nil {
			return err
		}
		if r.ID == 0 {
			return errors.New("not found")
		}
		//删除防火墙规则
		if err = s.Core.DeleteIPRule(r.Ip, r.Strategy); err != nil {
			return err
		}
		//重载防火墙
		if err = s.Core.Reload(); err != nil {
			return err
		}
		//删除数据库规则
		if err = r.Delete(global.PanelDB, &model.ConditionsT{}); err != nil {
			return err
		}
	}
	return nil
}

// GetIpRules 获取IP规则列表
func (s *SystemFirewallService) GetIpRules(query string, page, limit int) ([]*model.FirewallRuleIp, int64, error) {
	//获取列表
	rules, total, err := (&model.FirewallRuleIp{}).List(global.PanelDB, &model.ConditionsT{"ORDER": "create_time DESC"}, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

// UpdateIPRule 更新IP规则
func (s *SystemFirewallService) UpdateIPRule(id int64, rule *request.UpdateIPRuleR) error {
	//从数据库中获取规则
	r, err := (&model.FirewallRuleIp{ID: id}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if r.ID == 0 {
		return errors.New("not found")
	}
	//删除防火墙规则
	if err = s.Core.DeleteIPRule(r.Ip, r.Strategy); err != nil {
		return err
	}
	//重载防火墙
	if err = s.Core.Reload(); err != nil {
		return err
	}
	//更新数据库
	r.Ip = rule.SourceIp
	r.Strategy = rule.Strategy
	r.Ps = rule.Ps
	if err = r.Update(global.PanelDB); err != nil {
		return err
	}
	//添加防火墙规则
	if err = s.Core.CreateIPRule(r.Ip, r.Strategy); err != nil {
		return err
	}
	//重载防火墙
	if err = s.Core.Reload(); err != nil {
		return err
	}
	return nil
}

// BatchCreateForwardRule 批量创建转发规则
func (s *SystemFirewallService) BatchCreateForwardRule(rules []*request.CreateForwardRuleR) error {
	for _, rule := range rules {
		//检查数据库中源端口是否已存在
		_, total, err := (&model.FirewallRuleForward{}).List(global.PanelDB, &model.ConditionsT{
			"source_port = ?": rule.SourcePort,
			"protocol = ?":    rule.Protocol,
			"ORDER":           "create_time DESC",
		}, 0, 0)
		if err != nil {
			return err
		}
		if total > 0 {
			return errors.New(helper.Message("firewalld.ForwardRuleAlreadyExists"))
		}
		if err := s.Core.CreateForwardRule(rule.TargetIp, rule.Protocol, rule.SourcePort, rule.TargetPort); err != nil {
			return err
		}
		//重载防火墙
		if err := s.Core.Reload(); err != nil {
			return err
		}
		//保存到数据库
		insertData := &model.FirewallRuleForward{
			SourcePort: rule.SourcePort,
			TargetIp:   rule.TargetIp,
			TargetPort: rule.TargetPort,
			Protocol:   rule.Protocol,
			Ps:         rule.Ps,
		}
		_, err = (insertData).Create(global.PanelDB)
		if err != nil {
			return err
		}
	}
	return nil
}

// BatchDeleteForwardRule 批量删除转发规则
func (s *SystemFirewallService) BatchDeleteForwardRule(ids []int64) error {
	for _, id := range ids {
		//从数据库中获取规则
		r, err := (&model.FirewallRuleForward{ID: id}).Get(global.PanelDB)
		if err != nil {
			return err
		}
		if r.ID == 0 {
			return errors.New("not found")
		}
		//删除防火墙规则
		if err = s.Core.DeleteForwardRule(r.TargetIp, r.Protocol, r.SourcePort, r.TargetPort); err != nil {
			return err
		}
		//重载防火墙
		if err = s.Core.Reload(); err != nil {
			return err
		}
		//删除数据库规则
		if err = r.Delete(global.PanelDB, &model.ConditionsT{}); err != nil {
			return err
		}
	}
	return nil
}

// GetForwardRules 获取转发规则列表
func (s *SystemFirewallService) GetForwardRules(query string, page, limit int) ([]*model.FirewallRuleForward, int64, error) {
	//获取列表
	rules, total, err := (&model.FirewallRuleForward{}).List(global.PanelDB, &model.ConditionsT{"ORDER": "create_time DESC"}, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

// UpdateForwardRule 更新转发规则
func (s *SystemFirewallService) UpdateForwardRule(id int64, rule *request.UpdateForwardRuleR) error {
	//从数据库中获取规则
	r, err := (&model.FirewallRuleForward{ID: id}).Get(global.PanelDB)
	if err != nil {
		return err
	}
	if r.ID == 0 {
		return errors.New("not found")
	}
	//删除防火墙规则
	if err = s.Core.DeleteForwardRule(r.TargetIp, r.Protocol, r.SourcePort, r.TargetPort); err != nil {
		return err
	}
	//重载防火墙
	if err = s.Core.Reload(); err != nil {
		return err
	}
	//更新数据库
	r.SourcePort = rule.SourcePort
	r.TargetIp = rule.TargetIp
	r.TargetPort = rule.TargetPort
	r.Protocol = rule.Protocol
	r.Ps = rule.Ps
	if err = r.Update(global.PanelDB); err != nil {
		return err
	}
	//添加防火墙规则
	if err = s.Core.CreateForwardRule(r.TargetIp, r.Protocol, r.SourcePort, r.TargetPort); err != nil {
		return err
	}
	//重载防火墙
	if err = s.Core.Reload(); err != nil {
		return err
	}
	return nil
}
