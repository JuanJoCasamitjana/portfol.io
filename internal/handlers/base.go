package handlers

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/JuanJoCasamitjana/portfol.io/internal/utils"
	"github.com/labstack/echo/v4"
)

var IsAccessRestricted = false

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

func RestraintAccessMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		locale := utils.GetLocale(c)
		data := map[string]any{
			"locale":  locale,
			"title":   utils.Translate(locale, "403_title"),
			"message": utils.Translate(locale, "403_message"),
		}
		user, err := GetUserOfSession(c)
		if IsAccessRestricted && err != nil {
			return c.Render(403, "403", data)
		}
		if IsAccessRestricted && user.Authority.Level < model.AUTH_ADMIN.Level {
			return c.Render(403, "403", data)
		}
		return next(c)
	}
}
