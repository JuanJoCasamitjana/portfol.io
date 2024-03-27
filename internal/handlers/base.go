package handlers

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/JuanJoCasamitjana/portfol.io/internal/utils"
	"github.com/labstack/echo/v4"
)

func RenderIndex(c echo.Context) error {
	locale := utils.GetLocale(c)
	isAuthenticated := false
	user, err := GetUserOfSession(c)
	if err == nil {
		isAuthenticated = true
	}
	isModerator := user.Authority.Level == model.AUTH_MODERATOR.Level
	isAdmin := user.Authority.Level == model.AUTH_ADMIN.Level
	data := map[string]any{
		"title":           "Portfol.io",
		"locale":          locale,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"IsActive":        user.Active,
	}
	return c.Render(200, "index", data)
}

func RenderNavbar(c echo.Context) error {
	locale := utils.GetLocale(c)
	isAuthenticated := false
	user, err := GetUserOfSession(c)
	if err == nil {
		isAuthenticated = true
	}
	isModerator := user.Authority.Level == model.AUTH_MODERATOR.Level
	isAdmin := user.Authority.Level == model.AUTH_ADMIN.Level
	data := map[string]any{
		"locale":          locale,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"IsActive":        user.Active,
	}
	return c.Render(200, "navbar", data)
}

func SendFavicon(c echo.Context) error {
	return c.File("web/static/favicon.ico")
}
