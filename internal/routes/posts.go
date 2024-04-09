package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func setUpPostsRoutes(e *echo.Echo) {
	e.GET("/main", handlers.GetPostsMain)
	e.GET("/posts", handlers.GetPostsPaginated)
	e.GET("/posts/all", handlers.GetPostsSearch)
	e.GET("/posts/all/search", handlers.PostsSearchPaginated)
	e.GET("/posts/articles", handlers.GetArticleSearch)
	e.GET("/posts/articles/search", handlers.ArticleSearchPaginated)
	e.GET("/posts/galleries", handlers.GetGallerySearch)
	e.GET("/posts/galleries/search", handlers.GallerySearchPaginated)
	e.GET("/posts/moderation/tab", handlers.GetPostsModerationTab)
	e.GET("/posts/moderation", handlers.GetAllPostsForModeration)
	e.DELETE("/posts/moderation/:id", handlers.DeletePostModerators)
	//Articles
	e.GET("/article/:id", handlers.GetArticleByID)
	e.GET("/article/create", handlers.CreateArticleForm)
	e.POST("/article/create", handlers.CreateArticle)
	e.POST("/article/publish", handlers.CreateAndPublishArticle)
	e.GET("/article/mine", handlers.GetMyArticles)
	e.POST("/article/:id/tags/:tag", handlers.AddTagToArticle)
	e.GET("/article/:id/tags", handlers.GetTagsOfArticle)
	e.GET("/article/edit/:id", handlers.EditArticleForm)
	e.POST("/article/edit/:id", handlers.EditArticle)
	e.POST("/article/publish/:id", handlers.PublishArticle)
	e.DELETE("/article/delete/:id", handlers.DeleteArticle)
	e.GET("/article/tag/:name", handlers.ArticlesByTagPaginated)
	//Galleries
	e.GET("/gallery/create", handlers.CreateGallery)
	e.POST("/gallery/:id/images", handlers.AddImageToGallery)
	e.GET("/gallery/:id/images", handlers.GetImagesOfGallery)
	e.POST("/gallery/:id/title", handlers.ChangeTitleOfGallery)
	e.GET("/gallery/:id/title", handlers.GetChangeTitleOfGallery)
	e.POST("/gallery/:id/publish", handlers.PublishGallery)
	e.GET("/gallery/:id", handlers.GetGalleryByID)
	e.DELETE("/image/:id", handlers.DeleteImage)
	e.GET("/gallery/image-upload-form/:id", handlers.GetImageUploadForm)
	e.GET("/gallery/mine", handlers.GetMyGalleries)
	e.GET("gallery/:id/tags", handlers.GetTagsOfGallery)
	e.POST("/gallery/:id/tags/:tag", handlers.AddTagToGallery)
	e.GET("/gallery/tag/:name", handlers.GalleriesByTagPaginated)
	e.DELETE("/gallery/delete/:id", handlers.DeleteGallery)
	e.GET("/gallery/edit/:id", handlers.EditGalleryForm)
	e.POST("/gallery/edit/:id", handlers.EditGallery)
	//Tags
	e.POST("/tag/create", handlers.CreateTag)
	e.GET("/tag/create", handlers.CreateTagForm)
	e.GET("/tag/find", handlers.FindTags)

}
