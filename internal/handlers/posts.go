package handlers

import (
	"fmt"
	"html/template"
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
		next_page_loader = fmt.Sprintf(`<div hx-get="/posts?page=%d" hx=swap="outerHTML"
		hx-trigger="revealed" hx-replace-url="/posts" class="mb-3 p-3"></div>`, next_page)
	}
	data := map[string]any{
		"posts":     posts_content,
		"next_page": template.HTML(next_page_loader),
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
			}
		}
	}
	return posts_content
}

// Mostly about articles
func CreateArticleForm(c echo.Context) error {
	data := map[string]any{
		"locale": utils.GetLocale(c),
	}
	return c.Render(200, "article_form", data)
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
	if err != nil {
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
	if err != nil {
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
		next_page_loader = fmt.Sprintf(`<div hx-get="/my_articles?page=%d" hx=swap="outerHTML"
		hx-trigger="revealed" class="mb-3 p-3"></div>`, next_page)
	}
	data := map[string]any{
		"articles":  articles,
		"locale":    locale,
		"more":      more,
		"next_page": template.HTML(next_page_loader),
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

func GetArticleByID(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, _ := GetUserOfSession(c)
	id_str := c.Param("id")
	id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	article, err := database.FindArticleByID(id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	isAuthor := user.Username == article.Author
	data := map[string]any{
		"id":        article.ID,
		"title":     article.Title,
		"author":    article.Author,
		"createdAt": article.CreatedAt.Format("2006-01-02 15:04:05"),
		"updatedAt": article.UpdatedAt.Format("2006-01-02 15:04:05"),
		"content":   template.HTML(article.Content),
		"locale":    locale,
		"isAuthor":  isAuthor,
	}
	return c.Render(200, "article", data)
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
		"text":  template.HTML(article.Content),
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
	err = database.DeleteArticle(&article)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	return c.Render(200, "success", nil)
}

func AddTagToArticle(c echo.Context) error {
	_, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	id_str := c.Param("id")
	article_id, err := strconv.ParseUint(id_str, 10, 64)
	if err != nil {
		return c.String(400, "Bad Request")
	}
	tag_str := c.Param("tag")
	tag, err := database.FindTagByName(tag_str)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	article, err := database.FindArticleByID(article_id)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	err = database.AddTagToArticle(&article, &tag)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	c.Response().Header().Set("HX-Trigger", "tags-reload")
	return c.String(200, "Tag added successfully!")
}

//Mostly about galleries and images

// Since its a collection of images it's better to create it first
func CreateGallery(c echo.Context) error {
	var gallery model.Gallery
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	gallery.Author = user.Username
	err = database.CreateGallery(&gallery)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	data := map[string]any{
		"locale": utils.GetLocale(c),
		"id":     gallery.ID,
	}
	return c.Render(200, "gallery_form", data)
}

func AddImageToGallery(c echo.Context) error {
	idstr := c.Param("id")
	gallery_id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(500, "error", nil)
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
	c.Response().Header().Set("HX-Trigger", "gallery-reload")
	data := map[string]any{
		"id":     gallery.ID,
		"locale": utils.GetLocale(c),
	}
	return c.Render(200, "upload_image", data)
}

func GetChangeTitleOfGallery(c echo.Context) error {
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

func GetGalleryByID(c echo.Context) error {
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
	images := convertImagesToDataMap(gallery.Images)
	data := map[string]any{
		"id":        gallery.ID,
		"title":     gallery.Title,
		"author":    gallery.Author,
		"createdAt": gallery.CreatedAt,
		"updatedAt": gallery.UpdatedAt,
		"images":    images,
		"locale":    locale,
	}
	return c.Render(200, "gallery", data)
}

func GetMyGalleries(c echo.Context) error {
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
		next_page_loader = fmt.Sprintf(`<div hx-get="/gallery/mine?page=%d" hx=swap="outerHTML"
		hx-trigger="revealed" class="mb-3 p-3"></div>`, next_page)
	}
	data := map[string]any{
		"locale":    locale,
		"galleries": galleries,
		"more":      more,
		"nextPage":  template.HTML(next_page_loader),
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

// Tags are used in articles, projects and galleries
func CreateTagForm(c echo.Context) error {
	_, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	data := map[string]any{
		"locale": utils.GetLocale(c),
	}
	return c.Render(200, "tag_form", data)
}

func CreateTag(c echo.Context) error {
	locale := utils.GetLocale(c)
	_, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	name := c.FormValue("name")
	_, err = database.FindTagByName(name)
	if err == nil {

		data := map[string]any{
			"locale": locale,
			"error":  utils.Translate("tag_already_exists", locale),
		}
		return c.Render(200, "tag_form", data)
	}
	tag := model.Tag{Name: name}
	err = database.CreateTag(&tag)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	return c.Render(200, "success", nil)
}

func FindTags(c echo.Context) error {
	locale := utils.GetLocale(c)
	_, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	query := c.FormValue("query")
	tags_db, err := database.FindTagLikeName(query)
	if err != nil {
		return c.String(500, "Internal Server Error")
	}
	tags := make([]map[string]any, len(tags_db))
	for i := range tags_db {
		tags[i] = map[string]any{
			"name": tags_db[i].Name,
		}
	}
	data := map[string]any{
		"tags":   tags,
		"locale": locale,
	}
	return c.Render(200, "tags", data)
}
