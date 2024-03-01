package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func SetUpPublicationRotes(app *echo.Echo) {
	app.GET("/posts", handlers.GetPostsPaginated)
	app.GET("/posts/:username", handlers.GetPostsByUserPaginated)
	articles := app.Group("/articles")
	articles.POST("", handlers.CreateArticle)
	articles.GET("/mine", handlers.FindAllMyArticles)
	articles.GET("/create", handlers.CreateArticleForm)
	articles.GET("/:id", handlers.FindArticleByID)
	articles.GET("/edit/:id", handlers.EditArticleForm)
	articles.POST("/edit/:id", handlers.EditArticle)
	articles.GET("", handlers.GetTenArticles)
	articles.DELETE("/:id", handlers.DeleteArticle)
	images := app.Group("/images")
	images.GET("/create", handlers.CreateImageForm)
	images.GET("/collection/create", handlers.GetImageCollectionCreationForm)
	images.POST("/collection/:id/add", handlers.AddImageToImageCollection)
	images.GET("/collection/:id/image-list", handlers.GetImagesOfCollection)
	images.GET("/collection/mine", handlers.GetMyImageCollections)
	images.GET("/collection/:id", handlers.GetImageCollectionById)
	images.POST("/collection/:id/publish", handlers.PublishImageCollection)
	images.GET("/collection/:id/show", handlers.ShowImageCollection)
	images.GET("/collection/:id/edit-title", handlers.GetImageCollectionTitleEdit)
	images.POST("/collection/:id/edit-title", handlers.ChangeTitleOfImageCollection)
	images.DELETE("/:id", handlers.RemoveImageFromCollection)
	images.DELETE("/collection/:id", handlers.DeleteFullImageCollection)

}
