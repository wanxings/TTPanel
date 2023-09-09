package conf

import (
	"fmt"
	"gorm.io/gorm/logger"
	"strings"
)

type Sqlite struct {
	PanelPath string `mapstructure:"panel_path" json:"panel_path" yaml:"panel_path"` // 面板数据库路径
	TTWafPath string `mapstructure:"ttwaf_path" json:"ttwaf_path" yaml:"ttwaf_path"` // tt_waf数据库路径
	LogLevel  string `mapstructure:"log_level" json:"log_level" yaml:"log_level"`    // 日志级别
}

func (s *Sqlite) PanelDsn() string {
	return fmt.Sprintf("%s",
		s.PanelPath,
	)
}

func (s *Sqlite) TTWafDsn() string {
	return fmt.Sprintf("%s",
		s.TTWafPath,
	)
}

func (s *Sqlite) GetLogLevel() logger.LogLevel {
	switch strings.ToLower(s.LogLevel) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Error
	}
}
