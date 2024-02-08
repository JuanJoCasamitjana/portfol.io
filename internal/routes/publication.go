package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func SetUpPublicationRotes(app *echo.Echo) {
	app.GET("/posts", handlers.GetPostsPaginated)
	articles := app.Group("/articles")
	articles.POST("", handlers.CreateArticle)
	articles.GET("/mine", handlers.FindAllMyArticles)
	articles.GET("/create", handlers.CreateArticleForm)
	articles.GET("/:id", handlers.FindArticleByID)
	articles.GET("/edit/:id", handlers.EditArticleForm)
	articles.GET("", handlers.GetTenArticles)
	images := app.Group("/images")
	images.POST("", handlers.CreateImage)
	images.GET("/mine", handlers.FindAllMyImages)
	images.GET("/create", handlers.CreateImageForm)
	images.GET("/:id", handlers.FindImageByID)
}
