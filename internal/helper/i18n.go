package helper

import (
	"TTPanel/internal/global"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"strings"
)

func Message(ID string) string {
	ler := i18n.NewLocalizer(global.I18n, global.Config.System.Language)
	content, err := ler.Localize(&i18n.LocalizeConfig{
		MessageID: ID,
	})
	if err != nil {
		return ID
	}
	content = strings.ReplaceAll(content, "<no value>", "#")
	if content == "" {
		return ID
	} else {
		return content
	}
}
func MessageWithMap(ID string, TemplateData any) string {
	ler := i18n.NewLocalizer(global.I18n, global.Config.System.Language)
	content, err := ler.Localize(&i18n.LocalizeConfig{
		MessageID:    ID,
		TemplateData: TemplateData,
	})
	if err != nil {
		return fmt.Sprintf("%s %v", ID, TemplateData)
	}
	content = strings.ReplaceAll(content, "<no value>", "#")
	if content == "" {
		return ID
	} else {
		return content
	}

}
