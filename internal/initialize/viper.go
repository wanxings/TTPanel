package initialize

import (
	"TTPanel/internal/global"
	"fmt"
	"github.com/spf13/viper"
)

func InitViper() *viper.Viper {
	vp := viper.New()
	vp.SetConfigFile("/www/panel/config/config.yaml")
	//vp.SetConfigName("config")
	//vp.AddConfigPath(".")
	//vp.AddConfigPath("config/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config YamlFile: %s \n", err))
	}
	if err = vp.Unmarshal(&global.Config); err != nil {
		panic(err.Error())
	}
	return vp
}
