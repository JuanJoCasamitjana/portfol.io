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
	app.GET("/logout", handlers.Logout)
	users := app.Group("/user")
	users.GET("/edit", handlers.GetEditUser)
	users.GET("/password", handlers.GetPasswordEdit)
	users.PUT("/password", handlers.ChangePassword)
	users.GET("/profile", handlers.GetMyProfile)
	users.GET("/profile/:username", handlers.GetUserByUsername)

}
