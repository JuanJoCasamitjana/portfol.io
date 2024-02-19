package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	article, err := database.GetArticleById(id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	isAuthor := false
	if sess.Values["user_id"] != nil {
		user_id := sess.Values["user_id"].(uint64)
		user, err := database.GetUserByID(user_id)

		if err != nil {
			return c.Render(http.StatusInternalServerError, "error.html", err.Error())
		}
		isAuthor = article.Author == user.Username
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
	data.Others = map[string]interface{}{"id": article.ID}
	return c.Render(http.StatusOK, "edit_article", data)
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
	file, err := c.FormFile("image")
	data.Errors = make(map[string]string)
	data.FormValues = make(map[string]string)
	if err != nil {
		data.ErrorsExist = true
		data.Errors["image"] = "File cannot be empty. "
	}
	if data.ErrorsExist {
		data.FormValues = map[string]string{"footer": footer}
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
		case "image_collections":
			post, err := database.GetImageCollectionByIdIfPublished(post.OwnerID)
			if err == gorm.ErrRecordNotFound {
				continue
			}
			if err != nil {
				println(err.Error())
				return c.Render(http.StatusInternalServerError, "error.html", err.Error())
			}
			dataItem = map[string]interface{}{
				"post_type": "image_collections",
				"id":        post.ID,
				"author":    post.Author,
				"createdAt": post.CreatedAt.Format("02/01/2006 15:04:05"),
				"url":       post.Images[0].FilePath,
				"title":     post.Title,
				"amount":    len(post.Images),
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

func GetImageCollectionCreationForm(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	var imgColl model.ImageCollection
	imgColl.Author = user.Username
	if err = database.CreateImageCollection(&imgColl); err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	data := make(map[string]interface{})
	data["id"] = imgColl.ID
	return c.Render(http.StatusOK, "create_image_collection", data)
}

func AddImageToImageCollection(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	idstr := c.Param("id")
	collID, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusBadRequest, "error", nil)
	}
	imgColl, err := database.GetImageCollectionById(collID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	if imgColl.Author != user.Username {
		return c.Render(http.StatusUnauthorized, "error", nil)
	}
	var image model.Image
	title := c.FormValue("title")
	footer := c.FormValue("footer")
	file, err := c.FormFile("image")
	if err != nil {
		return c.Render(http.StatusUnauthorized, "error", nil)
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
	imgColl.Title = title
	err = database.UpdateImageCollection(&imgColl)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}

	image.Author = user.Username
	image.Footer = footer
	image.FilePath = "/static/images/" + fileName

	err = database.CreateImage(&image)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}

	err = database.AddImageToImageCollection(&imgColl, &image)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data := make(map[string]interface{})
	data["id"] = imgColl.ID
	data["value_title"] = title

	c.Response().Header().Set("HX-Trigger", "collection-updated")
	return c.Render(http.StatusCreated, "add_image_to_collection", data)
}

func GetImagesOfCollection(c echo.Context) error {
	idsrt := c.Param("id")
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	collID, err := strconv.ParseUint(idsrt, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	imgColl, err := database.GetImageCollectionById(collID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data := make([]map[string]interface{}, len(imgColl.Images))
	for i, img := range imgColl.Images {
		data[i] = map[string]interface{}{
			"url":      img.FilePath,
			"author":   img.Author,
			"footer":   img.Footer,
			"id":       img.ID,
			"isAuthor": img.Author == user.Username,
		}
	}
	return c.Render(http.StatusOK, "images", data)
}

func PublishImageCollection(c echo.Context) error {
	log.Println("Publishing collection")
	idstr := c.Param("id")
	sess, err := session.Get("session", c)
	if err != nil {
		log.Println("Error getting session")
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	collID, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		log.Println("Error parsing id")
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		log.Println("Error getting user")
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	imgColl, err := database.GetImageCollectionById(collID)
	if err != nil {
		log.Println("Error getting collection")
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	if imgColl.Author != user.Username {
		log.Println("User is not the author")
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	imgColl.Published = true
	if err := database.UpdateImageCollection(&imgColl); err != nil {
		log.Println("Error updating collection")
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	log.Println("Collection published")
	return c.Render(http.StatusOK, "success", nil)
}

func GetMyImageCollections(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	imgColls, err := database.GetAllImageCollectionsOfUserWithImages(user.Username)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	data := make([]map[string]interface{}, len(imgColls))
	for i, coll := range imgColls {
		url := ""
		if len(coll.Images) > 0 {
			url = coll.Images[0].FilePath
		}
		data[i] = map[string]interface{}{
			"id":        coll.ID,
			"title":     coll.Title,
			"createdAt": coll.CreatedAt.Format("02/01/2006 15:04:05"),
			"published": coll.Published,
			"url":       url,
			"author":    coll.Author,
			"amount":    len(coll.Images),
		}
	}
	return c.Render(http.StatusOK, "image_collections", data)
}

func GetImageCollectionById(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	imgColl, err := database.GetImageCollectionById(id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data := map[string]interface{}{
		"id":          imgColl.ID,
		"value_title": imgColl.Title,
	}
	return c.Render(http.StatusOK, "create_image_collection", data)
}

func ShowImageCollection(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	imgColl, err := database.GetImageCollectionByIdIfPublished(id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data_images := make([]map[string]interface{}, len(imgColl.Images))
	for i, img := range imgColl.Images {
		data_images[i] = map[string]interface{}{
			"url":    img.FilePath,
			"author": img.Author,
			"footer": img.Footer,
			"id":     img.ID,
		}
	}
	data := map[string]interface{}{
		"id":     imgColl.ID,
		"title":  imgColl.Title,
		"author": imgColl.Author,
		"images": data_images,
	}
	return c.Render(http.StatusOK, "collection_details", data)
}

func GetPostsByUserPaginated(c echo.Context) error {
	username := c.Param("username")
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	posts, err := database.GetPostsOfUser(username, page, 12)
	if err != nil {
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
		case "image_collections":
			post, err := database.GetImageCollectionByIdIfPublished(post.OwnerID)
			if err == gorm.ErrRecordNotFound {
				continue
			}
			if err != nil {
				return c.Render(http.StatusInternalServerError, "error.html", err.Error())
			}
			dataItem = map[string]interface{}{
				"post_type": "image_collections",
				"id":        post.ID,
				"author":    post.Author,
				"createdAt": post.CreatedAt.Format("02/01/2006 15:04:05"),
				"url":       post.Images[0].FilePath,
				"title":     post.Title,
				"amount":    len(post.Images),
			}
		}
		post_data[i] = dataItem
	}
	data := make(map[string]interface{})
	data["user"] = username
	data["posts"] = post_data
	data["more"] = len(posts) == 12
	data["next"] = page + 1
	return c.Render(http.StatusOK, "posts", data)
}

func GetImageCollectionTitleEdit(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	imgColl, err := database.GetImageCollectionById(id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data := map[string]interface{}{
		"id":          imgColl.ID,
		"value_title": imgColl.Title,
	}
	return c.Render(http.StatusOK, "change_title_of_collection", data)
}

func ChangeTitleOfImageCollection(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	title := c.FormValue("title")
	imgColl, err := database.GetImageCollectionById(id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	imgColl.Title = title
	if err := database.UpdateImageCollection(&imgColl); err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	data := map[string]interface{}{
		"id":          imgColl.ID,
		"value_title": imgColl.Title,
	}
	return c.Render(http.StatusOK, "title_of_collection", data)
}

func RemoveImageFromCollection(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	idstr := c.Param("id")
	imgID, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	img, err := database.GetImageById(imgID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	if img.Author != user.Username {
		return c.Render(http.StatusUnauthorized, "error.html", nil)
	}
	err = database.DeleteImage(&img)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	c.Response().Header().Set("HX-Trigger", "collection-updated")
	data := map[string]string{"message": "Image deleted successfully."}
	return c.JSON(http.StatusOK, data)
}

func DeleteFullImageCollection(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	idstr := c.Param("id")
	collID, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	imgColl, err := database.GetImageCollectionById(collID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	if imgColl.Author != user.Username {
		return c.Render(http.StatusUnauthorized, "error", nil)
	}
	for _, img := range imgColl.Images {
		err := database.DeleteImage(&img)
		if err != nil {
			return c.Render(http.StatusInternalServerError, "error", nil)
		}
	}
	err = database.DeleteImageCollection(&imgColl)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	return c.Render(http.StatusOK, "success", nil)
}

func DeleteArticle(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	idstr := c.Param("id")
	articleID, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	article, err := database.GetArticleById(articleID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	if article.Author != user.Username {
		return c.Render(http.StatusUnauthorized, "error", nil)
	}
	err = database.DeleteArticle(&article)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	data := map[string]string{"message": "Article deleted successfully."}
	return c.Render(http.StatusOK, "success", data)
}

func EditArticle(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	var data Data
	var article model.Article
	title := c.FormValue("title")
	text := c.FormValue("text")
	idstr := c.Param("id")
	id, err := strconv.ParseUint(idstr, 10, 64)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	article, err = database.GetArticleById(id)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}
	if article.Author != user.Username {
		return c.Render(http.StatusForbidden, "error.html", "You are not the author of this article")
	}
	article.Title = title
	article.Text = text
	err = database.UpdateArticle(&article)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", err.Error())
	}

	return c.Render(http.StatusOK, "success", data)
}
