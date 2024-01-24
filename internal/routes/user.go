package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func SetUpUserRoutes(app *echo.Echo) {
	app.GET("/register", handlers.RegisterUser)
	app.POST("/register", handlers.CreateUser)
	app.GET("/login", handlers.GetLogin)
	app.POST("/login", handlers.Login)
	//users := app.Group("/user")
}
