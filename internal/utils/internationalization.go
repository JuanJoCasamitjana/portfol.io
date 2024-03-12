package utils

import (
	"github.com/eduardolat/goeasyi18n"
)

var i18n *goeasyi18n.I18n

func init() {
	i18n = goeasyi18n.NewI18n(goeasyi18n.Config{
		FallbackLanguageName:    "en",
		DisableConsistencyCheck: false,
	})
	enTranslations, err := goeasyi18n.LoadFromJsonFiles("./web/locale/en/*.json")
	if err != nil {
		panic(err)
	}
	esTranslations, err := goeasyi18n.LoadFromJsonFiles("./web/locale/es/*.json")
	if err != nil {
		panic(err)
	}
	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)
}

func Translate(lang string, key string) string {
	if i18n.HasLanguage(lang) {
		return i18n.Translate(lang, key)
	}
	return i18n.Translate("en", key)
}
