package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func SetUpRoutes(app *echo.Echo) {
	app.GET("/favicon.ico", handlers.GetFavicon)
	app.GET("/", handlers.GetIndex)
	app.GET("/body", handlers.GetBody)
	app.GET("/navbar", handlers.GetNavbar)
	SetUpUserRoutes(app)
}
