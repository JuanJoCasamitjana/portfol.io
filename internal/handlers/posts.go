package handlers

import (
	"fmt"
	"html/template"
	"sort"
	"strconv"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/JuanJoCasamitjana/portfol.io/internal/utils"
	"github.com/labstack/echo/v4"
)

func GetPostsPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	posts, err := database.FindPostsPaginated(page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	posts_content := convertPostsToDataMap(posts)
	next_page := page + 1
	more := len(posts) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/posts?page=%d", next_page)
	}
	data := map[string]any{
		"posts":     posts_content,
		"next_page": template.HTML(next_page_loader), //skipcq  GSC-G203
		"more":      more,
		"locale":    locale,
	}
	return c.Render(200, "posts", data)
}

func convertPostsToDataMap(posts []model.Post) []map[string]interface{} {
	posts_content := make([]map[string]interface{}, len(posts))
	for i := range posts {
		switch posts[i].OwnerType {
		case "article":
			article, err := database.FindArticleByID(posts[i].OwnerID)
			if err != nil {
				continue
			}
			posts_content[i] = map[string]any{
				"id":        article.ID,
				"title":     article.Title,
				"author":    article.Author,
				"createdAt": article.CreatedAt,
				"updatedAt": article.UpdatedAt,
				"post_type": "article",
				"post_id":   posts[i].ID,
				"published": article.Published,
			}
		case "project":
			project, err := database.FindProjectByID(posts[i].OwnerID)
			if err != nil {
				continue
			}
			posts_content[i] = map[string]any{
				"id":        project.ID,
				"title":     project.Title,
				"author":    project.Author,
				"createdAt": project.CreatedAt,
				"updatedAt": project.UpdatedAt,
				"post_type": "project",
				"post_id":   posts[i].ID,
				"published": project.Published,
			}
		case "gallery":
			gallery, err := database.FindGalleryByID(posts[i].OwnerID)
			if err != nil {
				continue
			}
			num_images := len(gallery.Images)
			url := ""
			if num_images > 0 {
				url = gallery.Images[0].ThumbURL
			}
			posts_content[i] = map[string]any{
				"id":        gallery.ID,
				"title":     gallery.Title,
				"author":    gallery.Author,
				"createdAt": gallery.CreatedAt,
				"updatedAt": gallery.UpdatedAt,
				"post_type": "gallery",
				"url":       url,
				"amount":    num_images,
				"post_id":   posts[i].ID,
				"published": gallery.Published,
			}
		}
	}
	return posts_content
}

// Mostly about articles
func CreateArticleFormPart(c echo.Context) error {
	data := map[string]any{
		"locale": utils.GetLocale(c),
	}
	return c.Render(200, "article_form", data)
}

func CreateArticleFormFull(c echo.Context) error {
	user, err := GetUserOfSession(c)
	isAuthenticated := err == nil
	isModerator := isAuthenticated && user.Authority.Level == model.AUTH_MODERATOR.Level
	isAdmin := isAuthenticated && user.Authority.Level == model.AUTH_ADMIN.Level
	data := map[string]any{
		"locale":          utils.GetLocale(c),
		"isActive":        user.Active,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"app_title":       "Portfol.io",
	}
	return c.Render(200, "article_form_full", data)
}

func CreateArticleForm(c echo.Context) error {
	which := c.QueryParam("which")
	if which == "part" {
		return CreateArticleFormPart(c)
	}
	return CreateArticleFormFull(c)
}

func CreateArticle(c echo.Context) error {
	var article model.Article
	title, text := c.FormValue("title"), c.FormValue("text")
	data := map[string]any{
		"locale": utils.GetLocale(c),
		"title":  title,
		"text":   text,
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.Render(200, "article_form", data)
	}
	processedHTML, err := processHTML(text)
	if err != nil {
		return c.Render(200, "article_form", data)
	}
	article.Title = title
	article.Content = processedHTML
	article.Author = user.Username
	err = database.CreateArticle(&article)
	if err != nil {
		return c.Render(200, "article_form", data)
	}
	return c.Render(200, "success", nil)
}

func CreateAndPublishArticle(c echo.Context) error {
	var article model.Article
	title, text := c.FormValue("title"), c.FormValue("text")
	data := map[string]any{
		"locale": utils.GetLocale(c),
		"title":  title,
		"text":   text,
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.Render(200, "article_form", data)
	}
	processedHTML, err := processHTML(text)
	if err != nil {
		return c.Render(200, "article_form", data)
	}
	article.Title = title
	article.Content = processedHTML
	article.Author = user.Username
	article.Published = true
	err = database.CreateArticle(&article)
	if err != nil {
		return c.Render(200, "article_form", data)
	}
	return c.Render(200, "success", nil)
}

func GetMyArticles(c echo.Context) error {
	which := c.QueryParam("which")
	if which == "part" {
		return GetMyArticlesPart(c)
	}
	return GetMyArticlesFull(c)
}

func GetMyArticlesFull(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	isAuthenticated := err == nil
	isModerator := isAuthenticated && user.Authority.Level == model.AUTH_MODERATOR.Level
	isAdmin := isAuthenticated && user.Authority.Level == model.AUTH_ADMIN.Level
	data := map[string]any{
		"locale":          locale,
		"isActive":        user.Active,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"app_title":       "Portfol.io",
	}
	return c.Render(200, "article_list_full", data)
}

func GetMyArticlesPart(c echo.Context) error {
	locale := utils.GetLocale(c)
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	articles_db, err := database.FindAllArticlesByAuthorPaginated(user.Username, page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	articles := convertArticlesToDataMap(articles_db)
	next_page := page + 1
	more := len(articles_db) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/article/mine?page=%d&which=part", next_page)
	}
	data := map[string]any{
		"isMine":    true,
		"articles":  articles,
		"locale":    locale,
		"more":      more,
		"next_page": template.HTML(next_page_loader), //skipcq  GSC-G203
	}
	return c.Render(200, "article_list", data)
}

func convertArticlesToDataMap(articles []model.Article) []map[string]interface{} {
	articles_content := make([]map[string]interface{}, len(articles))
	for i := range articles {
		articles_content[i] = map[string]any{
			"id":        articles[i].ID,
			"title":     articles[i].Title,
			"author":    articles[i].Author,
			"createdAt": articles[i].CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt": articles[i].UpdatedAt.Format("2006-01-02 15:04:05"),
			"published": articles[i].Published,
		}
	}
	return articles_content
}

func GetArticleByIDPart(c echo.Context) error {
	locale := utils.GetLocale(c)
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	article, err := database.FindArticleByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	user, err := GetUserOfSession(c)
	if !article.Published && (err != nil || user.Username != article.Author) {
		return c.String(401, "Unauthorized")
	}
	isAuthor := user.Username == article.Author
	data := map[string]any{
		"id":        article.ID,
		"title":     article.Title,
		"author":    article.Author,
		"createdAt": article.CreatedAt.Format("2006-01-02 15:04:05"),
		"updatedAt": article.UpdatedAt.Format("2006-01-02 15:04:05"),
		"content":   template.HTML(article.Content), //skipcq  GSC-G203
		"published": article.Published,
		"locale":    locale,
		"isAuthor":  isAuthor,
		"isActive":  user.Active,
	}
	return c.Render(200, "article", data)
}

func GetArticleByIDFull(c echo.Context) error {
	locale := utils.GetLocale(c)
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	article, err := database.FindArticleByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	user, err := GetUserOfSession(c)
	if !article.Published && (err != nil || user.Username != article.Author) {
		return c.String(401, "Unauthorized")
	}
	isAuthor := user.Username == article.Author
	isAuthenticated := err == nil
	isModerator := isAuthenticated && user.Authority.Level == model.AUTH_MODERATOR.Level
	isAdmin := isAuthenticated && user.Authority.Level == model.AUTH_ADMIN.Level
	data := map[string]any{
		"id":              article.ID,
		"title":           article.Title,
		"content":         template.HTML(article.Content), //skipcq  GSC-G203
		"author":          article.Author,
		"createdAt":       article.CreatedAt,
		"updatedAt":       article.UpdatedAt,
		"published":       article.Published,
		"locale":          locale,
		"isAuthor":        isAuthor,
		"isActive":        user.Active,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"app_title":       "Portfol.io",
	}
	return c.Render(200, "article_full", data)
}

func GetArticleByID(c echo.Context) error {
	which := c.QueryParam("which")
	if which == "part" {
		return GetArticleByIDPart(c)
	}
	return GetArticleByIDFull(c)
}

func EditArticleForm(c echo.Context) error {
	locale := utils.GetLocale(c)
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	article, err := database.FindArticleByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	formValues := map[string]any{
		"title": article.Title,
		"text":  template.HTML(article.Content), //skipcq  GSC-G203
	}
	data := map[string]any{
		"id":         article.ID,
		"formValues": formValues,
		"locale":     locale,
	}
	return c.Render(200, "article_form", data)
}

func EditArticle(c echo.Context) error {
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	title, text := c.FormValue("title"), c.FormValue("text")
	data := map[string]any{
		"title":  title,
		"text":   text,
		"locale": utils.GetLocale(c),
	}
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.Render(200, "article_form", data)
	}
	if !user.Active {
		return c.String(401, "Unauthorized")
	}
	processedHTML, err := processHTML(text)
	if err != nil {
		return c.Render(200, "article_form", data)
	}
	article, err := database.FindArticleByID(id)
	if err != nil {
		return c.Render(200, "article_form", data)
	}
	if article.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	article.Title = title
	article.Content = processedHTML
	err = database.UpdateArticle(&article)
	if err != nil {
		return c.Render(200, "article_form", data)
	}
	return c.Render(200, "success", nil)
}

func PublishArticle(c echo.Context) error {
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.String(401, "Unauthorized")
	}
	article, err := database.FindArticleByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	if article.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	article.Published = true
	err = database.UpdateArticle(&article)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	return c.Render(200, "success", nil)
}

func DeleteArticle(c echo.Context) error {
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	article, err := database.FindArticleByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	if article.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	err = database.RemoveAllVotesForArticle(&article)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	err = database.DeleteArticle(&article)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	return c.Render(200, "success", nil)
}

//Mostly about galleries and images

// Since its a collection of images it's better to create it first
func CreateGalleryPart(c echo.Context) error {
	var gallery model.Gallery
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.String(401, "Unauthorized")
	}
	gallery.Author = user.Username
	err = database.CreateGallery(&gallery)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	data := map[string]any{
		"locale":      utils.GetLocale(c),
		"id":          gallery.ID,
		"isZero":      true,
		"isLimit":     false,
		"isPublished": gallery.Published,
	}
	return c.Render(200, "gallery_form", data)
}

func CreateGalleryFull(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	isAuthenticated := err == nil
	isModerator := isAuthenticated && user.Authority.Level == model.AUTH_MODERATOR.Level
	isAdmin := isAuthenticated && user.Authority.Level == model.AUTH_ADMIN.Level
	data := map[string]any{
		"locale":          locale,
		"isActive":        user.Active,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"app_title":       "Portfol.io",
		"page_to_load":    "/gallery/create?which=part",
	}
	return c.Render(200, "full_page_load", data)
}

func CreateGallery(c echo.Context) error {
	which := c.QueryParam("which")
	if which == "part" {
		return CreateGalleryPart(c)
	}
	return CreateGalleryFull(c)
}

func AddImageToGallery(c echo.Context) error {
	idstr := c.Param("id")
	gallery_id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(500, "error", nil)
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.String(401, "Unauthorized")
	}
	gallery, err := database.FindGalleryByID(gallery_id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	if len(gallery.Images) >= 10 {
		return c.String(400, "Bad Request")
	}
	if gallery.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	var image model.Image
	file, err := c.FormFile("image")
	if err != nil {
		return c.String(400, "Bad Request")
	}
	file_bytes, err := convertFileToBytes(file)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	image.GalleryID = gallery_id
	url_map, err := uploadImageToImgbb(file_bytes)
	if err != nil {
		return c.Render(500, "error", nil)
	}
	image.ImageURL = url_map["image_url"]
	image.ThumbURL = url_map["thumb_url"]
	image.DeleteURL = url_map["delete_url"]
	image.Footer = c.FormValue("footer")
	image.Owner = user.Username
	err = database.CreateImage(&image)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	amount := len(gallery.Images) + 1
	isLimit := amount >= 10
	c.Response().Header().Set("HX-Trigger", "gallery-reload")
	data := map[string]any{
		"id":          gallery.ID,
		"locale":      utils.GetLocale(c),
		"isLimit":     isLimit,
		"isZero":      false,
		"isPublished": gallery.Published,
	}
	return c.Render(200, "upload_image", data)
}

func GetChangeTitleOfGallery(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.String(401, "Unauthorized")
	}
	idstr := c.Param("id")
	gallery_id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	gallery, err := database.FindGalleryByID(gallery_id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	data := map[string]any{
		"id":          gallery.ID,
		"locale":      utils.GetLocale(c),
		"value_title": gallery.Title,
	}
	return c.Render(200, "gallery_title_form", data)
}

func ChangeTitleOfGallery(c echo.Context) error {
	idstr := c.Param("id")
	gallery_id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.String(401, "Unauthorized")
	}
	gallery, err := database.FindGalleryByID(gallery_id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	if gallery.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	title := c.FormValue("title")
	gallery.Title = title
	err = database.UpdateGallery(&gallery)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	data := map[string]any{
		"id":          gallery.ID,
		"locale":      utils.GetLocale(c),
		"value_title": title,
	}
	return c.Render(200, "gallery_title", data)
}

func PublishGallery(c echo.Context) error {
	idstr := c.Param("id")
	gallery_id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.String(401, "Unauthorized")
	}
	gallery, err := database.FindGalleryByID(gallery_id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	if gallery.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	gallery.Published = true
	err = database.UpdateGallery(&gallery)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	return c.Render(200, "success", nil)
}

func DeleteImage(c echo.Context) error {
	idstr := c.Param("id")
	image_id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	image, err := database.FindImageByID(image_id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	if image.Owner != user.Username {
		return c.String(401, "Unauthorized")
	}
	err = database.DeleteImage(&image)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	c.Response().Header().Set("HX-Trigger", "gallery-reload")
	data := map[string]string{
		"message": "Image deleted successfully!",
	}
	return c.JSON(200, data)
}

func GetImageUploadForm(c echo.Context) error {
	locale := utils.GetLocale(c)
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	gallery, err := database.FindGalleryByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active || gallery.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	isLimit := len(gallery.Images) >= 10
	isZero := len(gallery.Images) == 0
	data := map[string]any{
		"id":          gallery.ID,
		"locale":      locale,
		"isLimit":     isLimit,
		"isZero":      isZero,
		"isPublished": gallery.Published,
	}
	return c.Render(200, "upload_image", data)
}

func GetImagesOfGallery(c echo.Context) error {
	idstr := c.Param("id")
	gallery_id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	gallery, err := database.FindGalleryByID(gallery_id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	user, _ := GetUserOfSession(c)
	images := convertImagesToDataMap(gallery.Images, "isAuthor", user.Username == gallery.Author)
	data := map[string]any{
		"id":     gallery.ID,
		"images": images,
		"locale": utils.GetLocale(c),
	}
	return c.Render(200, "images", data)
}

func convertImagesToDataMap(images []model.Image, optional_values ...any) []map[string]interface{} {
	//Optional are key value pairs to be added to every image
	values := make(map[string]any, len(optional_values)/2)
	if len(optional_values)%2 == 0 {
		for i := 0; i < len(optional_values); i += 2 {
			values[optional_values[i].(string)] = optional_values[i+1]
		}
	}
	images_content := make([]map[string]interface{}, len(images))
	for i := range images {
		images_content[i] = map[string]any{
			"id":        images[i].ID,
			"image_url": images[i].ImageURL,
			"thumb_url": images[i].ThumbURL,
			"footer":    images[i].Footer,
			"author":    images[i].Owner,
			"options":   values,
		}
	}
	return images_content
}

func GetGalleryByIDPart(c echo.Context) error {
	locale := utils.GetLocale(c)
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	gallery, err := database.FindGalleryByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	user, err := GetUserOfSession(c)
	if !gallery.Published && (err != nil || user.Username != gallery.Author) {
		return c.String(401, "Unauthorized")
	}
	isAuthor := user.Username == gallery.Author
	images := convertImagesToDataMap(gallery.Images)
	data := map[string]any{
		"id":        gallery.ID,
		"title":     gallery.Title,
		"author":    gallery.Author,
		"createdAt": gallery.CreatedAt,
		"updatedAt": gallery.UpdatedAt,
		"published": gallery.Published,
		"images":    images,
		"locale":    locale,
		"isAuthor":  isAuthor,
		"isActive":  user.Active,
	}
	return c.Render(200, "gallery", data)
}

func GetGalleryByIDFull(c echo.Context) error {
	locale := utils.GetLocale(c)
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	gallery, err := database.FindGalleryByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	user, err := GetUserOfSession(c)
	if !gallery.Published && (err != nil || user.Username != gallery.Author) {
		return c.String(401, "Unauthorized")
	}
	isAuthor := user.Username == gallery.Author
	isAuthenticated := err == nil
	isModerator := isAuthenticated && user.Authority.Level == model.AUTH_MODERATOR.Level
	isAdmin := isAuthenticated && user.Authority.Level == model.AUTH_ADMIN.Level
	images := convertImagesToDataMap(gallery.Images)
	data := map[string]any{
		"id":              gallery.ID,
		"title":           gallery.Title,
		"author":          gallery.Author,
		"createdAt":       gallery.CreatedAt,
		"updatedAt":       gallery.UpdatedAt,
		"published":       gallery.Published,
		"images":          images,
		"locale":          locale,
		"isAuthor":        isAuthor,
		"isActive":        user.Active,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"app_title":       "Portfol.io",
	}
	return c.Render(200, "gallery_full", data)
}

func GetGalleryByID(c echo.Context) error {
	which := c.QueryParam("which")
	if which == "part" {
		return GetGalleryByIDPart(c)
	}
	return GetGalleryByIDFull(c)
}

func GetMyGalleries(c echo.Context) error {
	which := c.QueryParam("which")
	if which == "part" {
		return GetMyGalleriesPart(c)
	}
	return GetMyGalleriesFull(c)

}

func GetMyGalleriesFull(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	isAuthenticated := err == nil
	isModerator := isAuthenticated && user.Authority.Level == model.AUTH_MODERATOR.Level
	isAdmin := isAuthenticated && user.Authority.Level == model.AUTH_ADMIN.Level
	data := map[string]any{
		"locale":          locale,
		"isActive":        user.Active,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"app_title":       "Portfol.io",
		"page_to_load":    "/gallery/mine?page=1&which=part",
	}
	return c.Render(200, "full_page_load", data)
}

func GetMyGalleriesPart(c echo.Context) error {
	locale := utils.GetLocale(c)
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	galleries_db, err := database.FindAllGalleriesByAuthorPaginated(user.Username, page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	galleries := convertGalleriesToDataMap(galleries_db, true)
	next_page := page + 1
	more := len(galleries_db) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/gallery/mine?page=%d", next_page)
	}
	data := map[string]any{
		"locale":    locale,
		"galleries": galleries,
		"more":      more,
		"nextPage":  template.HTML(next_page_loader), //skipcq  GSC-G203
	}
	return c.Render(200, "gallery_list", data)
}

func DeleteGallery(c echo.Context) error {
	idstr := c.Param("id")
	gallery_id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	gallery, err := database.FindGalleryByID(gallery_id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	if gallery.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	err = database.RemoveAllVotesForGallery(&gallery)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	err = database.DeleteGallery(&gallery)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	return c.Render(200, "success", nil)
}

func convertGalleriesToDataMap(galleries_db []model.Gallery, showBadge bool) []map[string]any {
	galleries := make([]map[string]any, len(galleries_db))
	for i := range galleries_db {
		amount := len(galleries_db[i].Images)
		url := ""
		if amount > 0 {
			url = galleries_db[i].Images[0].ThumbURL
		}
		galleries[i] = map[string]any{
			"id":        galleries_db[i].ID,
			"title":     galleries_db[i].Title,
			"author":    galleries_db[i].Author,
			"createdAt": galleries_db[i].CreatedAt,
			"updatedAt": galleries_db[i].UpdatedAt,
			"published": galleries_db[i].Published,
			"url":       url,
			"amount":    amount,
			"showBadge": showBadge,
		}
	}
	return galleries
}

func EditGalleryForm(c echo.Context) error {
	locale := utils.GetLocale(c)
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	gallery, err := database.FindGalleryByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	isLimit := len(gallery.Images) >= 10
	isZero := len(gallery.Images) == 0
	data := map[string]any{
		"id":          gallery.ID,
		"locale":      locale,
		"value_title": gallery.Title,
		"isLimit":     isLimit,
		"isZero":      isZero,
		"isPublished": gallery.Published,
	}
	return c.Render(200, "gallery_form", data)
}

func EditGallery(c echo.Context) error {
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	title := c.FormValue("title")
	data := map[string]any{
		"title":  title,
		"locale": utils.GetLocale(c),
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.Render(200, "gallery_form", data)
	}
	gallery, err := database.FindGalleryByID(id)
	if err != nil {
		return c.Render(200, "gallery_form", data)
	}
	if gallery.Author != user.Username {
		return c.String(401, "Unauthorized")
	}
	gallery.Title = title
	err = database.UpdateGallery(&gallery)
	if err != nil {
		return c.Render(200, "gallery_form", data)
	}
	return c.Render(200, "success", nil)
}

func GalleriesByTagPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	tagName := c.Param("name")
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	galleries_db, err := database.FindAllGalleriesByTagPaginated(tagName, page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	galleries := convertGalleriesToDataMap(galleries_db, false)
	next_page := page + 1
	more := len(galleries_db) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/gallery/tag/%s?page=%d", tagName, next_page)
	}
	data := map[string]any{
		"galleries": galleries,
		"locale":    locale,
		"more":      more,
		"next_page": template.HTML(next_page_loader), //skipcq  GSC-G203
		"tag":       tagName,
	}
	return c.Render(200, "gallery_list", data)
}

// Tags are used in articles, projects and galleries
func CreateTagForm(c echo.Context) error {
	postType := c.QueryParam("post-type")
	postIDstr := c.QueryParam("post-id")
	data := map[string]any{
		"locale":   utils.GetLocale(c),
		"postType": postType,
		"postID":   postIDstr,
	}
	return c.Render(200, "tag_form", data)
}

func CreateTag(c echo.Context) error {
	locale := utils.GetLocale(c)
	postType := c.QueryParam("post-type")
	postIDstr := c.QueryParam("post-id")
	name := c.FormValue("query")
	_, err := database.FindTagByName(name)
	if err == nil {

		data := map[string]any{
			"locale":   locale,
			"error":    utils.Translate("tag_already_exists", locale),
			"postType": postType,
			"postID":   postIDstr,
		}
		return c.Render(200, "tag_form", data)
	}
	tag := model.Tag{Name: name}
	err = database.CreateTag(&tag)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	data := map[string]any{
		"locale":   locale,
		"postType": postType,
		"postID":   postIDstr,
	}
	return c.Render(200, "tag_form", data)
}

func FindTags(c echo.Context) error {
	locale := utils.GetLocale(c)
	query := c.QueryParam("query")
	postType := c.QueryParam("post-type")
	postIDstr := c.QueryParam("post-id")
	isAdd := false
	if postType != "" {
		isAdd = true
	}
	tags_db, err := database.FindTagLikeName(query, 25)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	postID, err := strconv.ParseUint(postIDstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	tags := make([]map[string]any, len(tags_db))
	for i := range tags_db {
		tags[i] = map[string]any{
			"name":      tags_db[i].Name,
			"add":       isAdd,
			"post_type": postType,
			"post_id":   postID,
			"bgColor":   tags_db[i].ColorOfTag(),
		}
	}
	data := map[string]any{
		"tags":   tags,
		"locale": locale,
	}
	return c.Render(200, "tags", data)
}

/*
	 func filterTags(toBeFiltered, filter []model.Tag) []model.Tag {
		var filtered []model.Tag
		for i := range toBeFiltered {
			add := true
			for j := range filter {
				if toBeFiltered[i].Name == filter[j].Name {
					add = false
					break
				}
			}
			if add {
				filtered = append(filtered, toBeFiltered[i])
			}
		}
		return filtered
	}
*/
func ArticlesByTagPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	tagName := c.Param("name")
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	articles_db, err := database.FindAllArticlesByTagPaginated(tagName, page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	articles := convertArticlesToDataMap(articles_db)
	next_page := page + 1
	more := len(articles_db) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/article/tag/%s?page=%d", tagName, next_page)
	}
	data := map[string]any{
		"articles":  articles,
		"locale":    locale,
		"more":      more,
		"next_page": template.HTML(next_page_loader), //skipcq  GSC-G203
		"tag":       tagName,
	}
	return c.Render(200, "article_list", data)
}

// Implementing  content filering

func GetPostsMain(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "posts_main", data)
}

func PostsMainPage(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "posts_main", data)
}

func PostsSearchPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	query := c.QueryParam("query")
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		page = 1
	}
	posts_db, err := database.FindPostsByQueryPaginated(query, page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	posts_content := convertPostsToDataMap(posts_db)
	next_page := page + 1
	more := len(posts_db) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/posts/all/search?query=%s&page=%d", query, next_page)
	}
	data := map[string]any{
		"posts":    posts_content,
		"nextPage": template.HTML(next_page_loader), //skipcq  GSC-G203
		"more":     more,
		"locale":   locale,
	}
	return c.Render(200, "posts", data)
}

func GetPostsSearch(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "posts_search", data)
}

func ArticleSearchPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	query := c.QueryParam("query")
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		page = 1
	}
	articles_db, err := database.FindArticlesByQueryPaginated(query, page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	articles := convertArticlesToDataMap(articles_db)
	next_page := page + 1
	more := len(articles_db) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/posts/articles/search?query=%s&page=%d", query, next_page)
	}
	data := map[string]any{
		"articles": articles,
		"nextPage": template.HTML(next_page_loader), //skipcq  GSC-G203
		"more":     more,
		"locale":   locale,
	}
	return c.Render(200, "article_list", data)
}

func GetArticleSearch(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "article_search", data)
}

func GallerySearchPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	query := c.QueryParam("query")
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		page = 1
	}
	galleries_db, err := database.FindGalleriesByQueryPaginated(query, page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	galleries := convertGalleriesToDataMap(galleries_db, false)
	next_page := page + 1
	more := len(galleries_db) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/posts/galleries/search?query=%s&page=%d", query, next_page)
	}
	data := map[string]any{
		"galleries": galleries,
		"nextPage":  template.HTML(next_page_loader), //skipcq  GSC-G203
		"more":      more,
		"locale":    locale,
	}
	return c.Render(200, "galleries", data)
}

func GetGallerySearch(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "gallery_search", data)
}

func GetPostsModerationTab(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "posts_moderation", data)
}

func GetAllPostsForModeration(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active || user.Authority.Level < model.AUTH_MODERATOR.Level {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	query := c.QueryParam("query")
	postsDB, err := database.FindAllPostsByqueryPaginated(page, 12, query)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	posts := convertPostsToPostList(postsDB)
	more := len(postsDB) == 12
	nextPage := page + 1
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/posts/moderation?page=%d&query=%s", nextPage, query)
	}
	data := map[string]any{
		"locale": locale,
		"posts":  posts,
		"more":   more,
		"next":   template.HTML(nextPageLoader), //skipcq  GSC-G203
	}
	return c.Render(200, "posts_moderation_list", data)
}

func convertPostsToPostList(posts []model.Post) []map[string]any {
	var postsContent []map[string]any
	for i := range posts {
		postsContent = append(postsContent, map[string]any{
			"postID":    posts[i].ID,
			"author":    posts[i].Author,
			"title":     posts[i].Title,
			"createdAt": posts[i].CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt": posts[i].UpdatedAt.Format("2006-01-02 15:04:05"),
			"published": posts[i].Published,
			"type":      posts[i].OwnerType,
			"id":        posts[i].OwnerID,
		})
	}
	return postsContent
}

func DeletePostModerators(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active || user.Authority.Level < model.AUTH_MODERATOR.Level {
		return c.String(401, "Unauthorized")
	}
	postIDstr := c.Param("id")
	postID, err := strconv.ParseUint(postIDstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	err = database.DeletePostByID(postID)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	return c.String(200, "Post deleted successfully!")
}

func VoteTagForPost(c echo.Context) error {
	tagName := c.QueryParam("tag")
	ownerIDstr := c.QueryParam("postid")
	ownerType := c.QueryParam("posttype")
	var vote model.Vote
	ownerID, err := strconv.ParseUint(ownerIDstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}

	article := new(model.Article)
	gallery := new(model.Gallery)
	switch ownerType {
	case "article":
		articleDB, err := database.FindArticleByID(ownerID)
		if err != nil {
			return c.String(500, "Internal Server Error")
		}
		article = &articleDB
	case "gallery":
		galleryDB, err := database.FindGalleryByID(ownerID)
		if err != nil {
			return c.String(500, "Internal Server Error")
		}
		gallery = &galleryDB
	default:
		return c.String(400, "Bad Request")
	}
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
		return c.String(401, "Unauthorized")
	}
	tag, err := database.FindTagByName(tagName)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	alreadyVoted := database.VoteExistsForTagUserAndPost(tag.ID, user.Username, ownerID, ownerType)
	if alreadyVoted {
		return c.String(403, "Forbidden")
	}
	vote.TagID = tag.ID
	vote.Voter = user.Username
	switch ownerType {
	case "article":
		err = database.VoteTagForArticle(article, &vote)
		if err != nil {
			return c.String(500, "Internal Server Error")
		}
	case "gallery":
		err = database.VoteTagForGallery(gallery, &vote)
		if err != nil {
			return c.String(500, "Internal Server Error")
		}
	}
	c.Response().Header().Set("HX-Trigger", "votes-reload")
	return c.String(200, "Tag voted successfully!")
}

func FindVotesOfGallery(c echo.Context) error {
	galleryIDstr := c.Param("id")
	galleryID, err := strconv.ParseUint(galleryIDstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	gallery, err := database.FindGalleryByID(galleryID)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	votes := convertVotesToDataMap(gallery.Votes)
	data := map[string]any{
		"votes":  votes,
		"locale": utils.GetLocale(c),
	}
	return c.Render(200, "votes", data)
}

func FindVotesOfArticle(c echo.Context) error {
	articleIDstr := c.Param("id")
	articleID, err := strconv.ParseUint(articleIDstr, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	gallery, err := database.FindArticleByID(articleID)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	votes := convertVotesToDataMap(gallery.Votes)
	data := map[string]any{
		"votes":  votes,
		"locale": utils.GetLocale(c),
	}
	return c.Render(200, "votes", data)
}

func convertVotesToDataMap(votes []model.Vote) []map[string]any {
	tagByNumberOfVotes := make(map[string][]any)
	var votesContent []map[string]any
	for i := range votes {
		if _, ok := tagByNumberOfVotes[votes[i].Tag.Name]; !ok {
			tagByNumberOfVotes[votes[i].Tag.Name] = []any{1, votes[i].Tag.ColorOfTag()}
		} else {
			tagByNumberOfVotes[votes[i].Tag.Name][0] = tagByNumberOfVotes[votes[i].Tag.Name][0].(int) + 1
		}
	}
	for tagName, tagInfo := range tagByNumberOfVotes {
		votesContent = append(votesContent, map[string]any{
			"tag":   tagName,
			"votes": tagInfo[0],
			"color": tagInfo[1],
		})
	}
	sort.Slice(votesContent, func(i, j int) bool {
		return votesContent[i]["votes"].(int) > votesContent[j]["votes"].(int)
	})
	return votesContent
}

func FindPostsByTagPaginatedPart(c echo.Context) error {
	locale := utils.GetLocale(c)
	tagName := c.Param("name")
	page_str := c.QueryParam("page")
	page, err := strconv.Atoi(page_str)
	if err != nil {
		page = 1
	}
	posts_db, err := database.FindPaginatedPostsByTagOrderedByNumberOfVotes(tagName, page, 12)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	posts_content := convertPostsToDataMap(posts_db)
	next_page := page + 1
	more := len(posts_db) == 12
	next_page_loader := ""
	if more {
		next_page_loader = fmt.Sprintf("/posts/all/tag/%s?page=%d", tagName, next_page)
	}
	data := map[string]any{
		"posts":    posts_content,
		"nextPage": template.HTML(next_page_loader), //skipcq  GSC-G203
		"more":     more,
		"locale":   locale,
		"tag":      tagName,
	}
	return c.Render(200, "posts", data)
}

func FindPostsByTagPaginatedFull(c echo.Context) error {
	locale := utils.GetLocale(c)
	isAuthenticated := false
	user, err := GetUserOfSession(c)
	if err == nil {
		isAuthenticated = true
	}
	tagName := c.Param("name")
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	isModerator := IsModerator(c)
	isAdmin := IsAdmin(c)
	pageToLoad := fmt.Sprintf("/posts/all/tag/%s?which=part&page=%d", tagName, page)
	data := map[string]any{
		"app_title":       "Portfol.io",
		"locale":          locale,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"isActive":        user.Active,
		"page_to_load":    pageToLoad,
	}
	return c.Render(200, "full_page_load", data)
}

func FindPostsByTagPaginated(c echo.Context) error {
	which := c.QueryParam("which")
	if which == "part" {
		return FindPostsByTagPaginatedPart(c)
	}
	return FindPostsByTagPaginatedFull(c)
}
