package global

import (
	"TTPanel/internal/conf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"sync"
)

var (
	PanelDB *gorm.DB
	TTWafDB *gorm.DB
	Config  conf.Server
	Vp      *viper.Viper
	I18n    *i18n.Bundle
	Log     *logrus.Logger
	Version string

	GoCache *cache.Cache

	//BlackCache local_cache.Cache
	lock sync.RWMutex
)
