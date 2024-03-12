package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func SetUpRoutes(e *echo.Echo) {
	e.GET("/", handlers.RenderIndex)
	e.GET("/navbar", handlers.RenderNavbar)
	e.GET("/favicon.ico", handlers.SendFavicon)
	setUpUsersRoutes(e)
	setUpPostsRoutes(e)
}
