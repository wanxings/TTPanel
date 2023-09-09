package initialize

import (
	"TTPanel/internal/global"
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"strings"
)

func InitI18n() *i18n.Bundle {

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	i18nPath := global.Config.System.PanelPath + "/data/i18n"
	files, err := os.ReadDir(i18nPath)
	if err != nil {
		global.Log.Error(err)
		os.Exit(1)
	}
	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), ".json") {
			bundle.MustLoadMessageFile(i18nPath + "/" + file.Name())
		}
	}

	global.Log.Debugf("panel use language:%s", global.Config.System.Language)
	return bundle
}
