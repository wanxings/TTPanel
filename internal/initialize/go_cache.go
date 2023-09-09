package initialize

import (
	"TTPanel/internal/global"
	"TTPanel/pkg/util"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

func InitGoCache() *cache.Cache {
	goCache := cache.New(5*time.Minute, 10*time.Minute)
	cacheFilePath := fmt.Sprintf("%s/data/cache.gob", global.Config.System.PanelPath)
	if util.IsFile(cacheFilePath) {
		err := goCache.LoadFile(cacheFilePath)
		if err != nil {
			global.Log.Errorf("InitGoCache->goCache.LoadFile Error:%s", err.Error())
		}
	}
	return goCache
}
