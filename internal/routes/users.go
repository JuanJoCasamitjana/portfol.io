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
	profile := e.Group("/profile")
	profile.GET("/mine", handlers.GetMyProfile)
	profile.GET("/mine/edit", handlers.GetProfileEditForm)
	profile.GET("/mine/edit/password", handlers.ChangePasswordForm)
	profile.POST("/mine/edit/password", handlers.ChangePassword)
	profile.DELETE("/mine", handlers.DeleteProfile)
	profile.POST("/mine/edit", handlers.EditProfile)
	profile.GET("/:username", handlers.GetUserProfile)
	profile.GET("/:username/sections", handlers.GetUserSections)
	profile.GET("/:username/sections/:section", handlers.GetUserSectionPaginated)
	profile.DELETE("/:username/sections/:section", handlers.DeleteSection)
	profile.GET("/:username/edit/sections", handlers.GetMySectionsList)
	profile.GET("/:username/edit/sections/:section", handlers.GetSectionEdit)
	profile.GET("/:username/create/section", handlers.CreateNewSectionForm)
	profile.POST("/:username/create/section", handlers.CreateNewSection)
	profile.PUT("/:username/sections/:section/edit/:post", handlers.AddPostToSection)
	profile.DELETE("/:username/sections/:section/edit/:post", handlers.RemovePostFromSection)
	profile.GET("/:username/posts/in/:section", handlers.GetModificablePostsFromSectionPaginated)
	profile.GET("/:username/posts/not-in/:section", handlers.GetModificablePostsNotFromSectionPaginated)
}
