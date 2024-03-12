package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func setUpUsersRoutes(e *echo.Echo) {
	e.GET("/register", handlers.GetRegisterForm)
	e.POST("/register", handlers.Register)
	e.GET("/login", handlers.GetLoginForm)
	e.POST("/login", handlers.Login)
	e.GET("/logout", handlers.Logout)
}
