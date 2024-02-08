package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func FindAllMyArticles(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	articles, err := database.GetAllArticlesOfUser(user.Username)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data := make([]map[string]interface{}, len(articles))
	for i, article := range articles {
		data[i] = make(map[string]interface{})
		data[i]["id"] = article.ID
		data[i]["title"] = article.Title
		data[i]["author"] = article.Author
		data[i]["createdAt"] = article.CreatedAt.Format("02/01/2006 15:04:05")
	}
	return c.Render(http.StatusOK, "articles", data)
}

func FindArticleByID(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	user_id := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(user_id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	article, err := database.GetArticleById(id)
	isAuthor := article.Author == user.Username
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	var data Data
	data.PureHTML = template.HTML(article.Text)
	data.Others = map[string]interface{}{
		"title":     article.Title,
		"author":    article.Author,
		"createdAt": article.CreatedAt.Format("02/01/2006 15:04:05"),
		"isAuthor":  isAuthor,
		"id":        article.ID,
	}
	return c.Render(http.StatusOK, "article", data)
}

func EditArticleForm(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	user_id := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(user_id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	article, err := database.GetArticleById(id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	if article.Author != user.Username {
		return c.Render(http.StatusForbidden, "error.html", "You are not the author of this article")
	}
	var data Data
	data.FormValues = map[string]string{"title": article.Title, "text": article.Text}
	data.PureHTML = template.HTML(article.Text)
	return c.Render(http.StatusOK, "create_article", data)
}

func FindAllMyImages(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	images, err := database.GetAllImagesOfUser(user.Username)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data_images := make([]map[string]interface{}, len(images))
	for i, image := range images {
		data_images[i] = make(map[string]interface{})
		data_images[i]["id"] = image.ID
		data_images[i]["title"] = image.Title
		data_images[i]["author"] = image.Author
		data_images[i]["createdAt"] = image.CreatedAt
		data_images[i]["url"] = image.FilePath
	}
	return c.Render(http.StatusOK, "images", data_images)
}

func FindAllMyArticleCollections(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	collections, err := database.GetAllArticleCollectionsOfUser(user.Username)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	return c.Render(http.StatusOK, "collections.html", collections)
}

func FindAllMyImageCollections(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	sesions, err := database.GetAllImageCollectionsOfUser(user.Username)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	return c.Render(http.StatusOK, "sesions.html", sesions)
}

func CreateArticleForm(c echo.Context) error {
	return c.Render(http.StatusOK, "create_article", nil)
}

func CreateArticle(c echo.Context) error {
	var data Data
	var article model.Article
	title := c.FormValue("title")
	text := c.FormValue("text")
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data.FormValues = map[string]string{"title": title, "text": text}
	data.PureHTML = template.HTML(text)
	if title == "" {
		data.Errors = make(map[string]string)
		data.ErrorsExist = true
		data.Errors["title"] = "Title cannot be empty. "
		return c.Render(http.StatusOK, "create_article", data)
	}
	article.Author = user.Username
	article.Title = title
	article.Text = text
	err = database.CreateArticle(&article)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	articles, err := database.GetAllArticlesOfUser(user.Username)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data_articles := make([]map[string]interface{}, len(articles))
	for i, article := range articles {
		data_articles[i] = make(map[string]interface{})
		textLeght := len(article.Text)
		preview := article.Text
		if textLeght > 100 {
			preview = preview[:100]
		}
		data_articles[i]["id"] = article.ID
		data_articles[i]["title"] = article.Title
		data_articles[i]["author"] = article.Author
		data_articles[i]["preview"] = preview + "..."
		data_articles[i]["createdAt"] = article.CreatedAt
	}
	return c.Render(http.StatusOK, "articles", data_articles)
}

func CreateImageForm(c echo.Context) error {
	return c.Render(http.StatusOK, "create_image", nil)
}

func CreateImage(c echo.Context) error {
	var data Data
	var image model.Image
	footer := c.FormValue("footer")
	title := c.FormValue("title")
	file, err := c.FormFile("image")
	data.Errors = make(map[string]string)
	data.FormValues = make(map[string]string)
	if err != nil {
		data.ErrorsExist = true
		data.Errors["image"] = "File cannot be empty. "
	}
	if title == "" {
		data.ErrorsExist = true
		data.Errors["title"] = "Title cannot be empty. "
	}
	if data.ErrorsExist {
		data.FormValues = map[string]string{"title": title, "footer": footer}
		return c.Render(http.StatusOK, "create_image", data)
	}
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	fileName := fmt.Sprintf("%s_%s", user.Username, time.Now().Format("2006-01-02_15-04-05"))
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dstPath := "web/static/images/" + fileName
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	image.Author = user.Username
	image.Footer = footer
	image.Title = title
	image.FilePath = "static/images/" + fileName
	err = database.CreateImage(&image)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	images, err := database.GetAllImagesOfUser(user.Username)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data_images := make([]map[string]interface{}, len(images))
	for i, image := range images {
		data_images[i] = make(map[string]interface{})
		data_images[i]["id"] = image.ID
		data_images[i]["title"] = image.Title
		data_images[i]["author"] = image.Author
		data_images[i]["createdAt"] = image.CreatedAt
		data_images[i]["url"] = image.FilePath
	}
	return c.Render(http.StatusOK, "images", data_images)
}

func GetTenArticles(c echo.Context) error {
	indexstr := c.QueryParam("index")
	index, err := strconv.Atoi(indexstr)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	articles, err := database.GetPaginatedArticles(index, 10)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data := make(map[string]interface{})
	articles_dat := make([]map[string]interface{}, len(articles))
	for i, article := range articles {
		articles_dat[i] = make(map[string]interface{})
		articles_dat[i]["id"] = article.ID
		articles_dat[i]["title"] = article.Title
		articles_dat[i]["author"] = article.Author
		articles_dat[i]["createdAt"] = article.CreatedAt.Format("02/01/2006 15:04:05")
	}
	data["articles"] = articles_dat
	data["more"] = len(articles) == 10
	data["next"] = index + 1
	return c.Render(http.StatusOK, "main-app", data)
}

func FindImageByID(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	image, err := database.GetImageById(id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	return c.Render(http.StatusOK, "image", image)
}

func GetPostsPaginated(c echo.Context) error {
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		println(err.Error())
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	posts, err := database.GetPosts(page, 12)
	if err != nil {
		println(err.Error())
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	post_data := make([]map[string]interface{}, len(posts))
	for i, post := range posts {
		var dataItem map[string]interface{}
		switch post.OwnerType {
		case "articles":
			post, err := database.GetArticleById(post.OwnerID)
			if err != nil {
				println(err.Error())
				return c.Render(http.StatusInternalServerError, "error.html", err.Error())
			}
			dataItem = map[string]interface{}{
				"post_type": "article",
				"id":        post.ID,
				"title":     post.Title,
				"author":    post.Author,
				"createdAt": post.CreatedAt.Format("02/01/2006 15:04:05"),
			}
		case "images":
			post, err := database.GetImageById(post.OwnerID)
			if err != nil {
				println(err.Error())
				return c.Render(http.StatusInternalServerError, "error.html", err.Error())
			}
			dataItem = map[string]interface{}{
				"post_type": "image",
				"id":        post.ID,
				"title":     post.Title,
				"author":    post.Author,
				"createdAt": post.CreatedAt.Format("02/01/2006 15:04:05"),
				"url":       post.FilePath,
			}
		}
		post_data[i] = dataItem
	}
	data := make(map[string]interface{})
	data["posts"] = post_data
	data["more"] = len(posts) == 12
	data["next"] = page + 1
	return c.Render(http.StatusOK, "posts", data)
}
