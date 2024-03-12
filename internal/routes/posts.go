package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func setUpPostsRoutes(e *echo.Echo) {
	e.GET("/posts", handlers.GetPostsPaginated)
	e.GET("/article/:id", handlers.GetArticleByID)
	e.GET("/gallery/:id", handlers.GetGalleryByID)
	e.GET("/article/create", handlers.CreateArticleForm)
	e.POST("/article/create", handlers.CreateArticle)
	e.POST("/article/publish", handlers.CreateAndPublishArticle)
	e.GET("/article/mine", handlers.GetMyArticles)
	e.POST("/article/:id/add-tag/:tag", handlers.AddTagToArticle)
	e.GET("/article/edit/:id", handlers.EditArticleForm)
	e.POST("/article/edit/:id", handlers.EditArticle)
	e.POST("/article/publish/:id", handlers.PublishArticle)
	e.DELETE("/article/delete/:id", handlers.DeleteArticle)
	e.GET("/gallery/create", handlers.CreateGallery)
	e.POST("/gallery/:id/images", handlers.AddImageToGallery)
	e.GET("/gallery/:id/images", handlers.GetImagesOfGallery)
	e.POST("/gallery/:id/title", handlers.ChangeTitleOfGallery)
	e.GET("/gallery/:id/title", handlers.GetChangeTitleOfGallery)
	e.POST("/galery/:id/publish", handlers.PublishGallery)
	e.GET("/gallery/:id", handlers.GetGalleryByID)
	e.DELETE("/image/:id", handlers.DeleteImage)
	e.GET("/gallery/mine", handlers.GetMyGalleries)
}
