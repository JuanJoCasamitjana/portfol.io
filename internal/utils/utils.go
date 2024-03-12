package utils

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func GetLocale(c echo.Context) string {
	lang := c.Request().Header.Get("Accept-Language")
	if lang == "" {
		return "en"
	}
	lang = strings.Split(lang, "-")[0]
	return lang
}
