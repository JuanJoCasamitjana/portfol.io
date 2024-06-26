package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func SetUpRoutes(e *echo.Echo) {
	e.Use(handlers.RestraintAccessMiddleware)
	e.GET("/", handlers.RenderIndex)
	e.GET("/navbar", handlers.RenderNavbar)
	e.GET("/favicon.ico", handlers.SendFavicon)
	e.GET("/admin/shutdown", handlers.ShutdownServer)
	setUpUsersRoutes(e)
	setUpPostsRoutes(e)
	setUpReportsRoutes(e)
}
