package handlers

import (
	"html/template"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Data struct {
	Title           string
	Message         string
	IsAuthenticated bool
	UserInfo        map[string]string
	Errors          map[string]string
	FormValues      map[string]string
	Others          map[string]interface{}
	ErrorsExist     bool
	PureHTML        template.HTML
}

func GetIndex(c echo.Context) error {
	var data Data
	data.Title = "Portfol.io"
	userInfo := make(map[string]string)
	data.IsAuthenticated = false
	sess, _ := session.Get("session", c)
	if sess.Values["user_id"] != nil {
		data.IsAuthenticated = true
		user, err := database.GetUserByID(sess.Values["user_id"].(uint64))
		if err != nil {
			data.IsAuthenticated = false
			return c.Render(c.Response().Status, "navbar.html", data)
		}
		userInfo["username"] = user.Username
		if user.Firstname != nil && user.Lastname != nil {
			userInfo["fullname"] = *user.Firstname + " " + *user.Lastname
		}
	}
	data.UserInfo = userInfo
	index_data := make(map[string]interface{})
	index_data["navbar"] = data
	articles, err := database.GetPaginatedArticles(1, 10)
	if err != nil {
		return c.Render(c.Response().Status, "index.html", data)
	}
	main_app_dat := make(map[string]interface{})
	articles_dat := make([]map[string]interface{}, len(articles))
	for i, article := range articles {
		articles_dat[i] = make(map[string]interface{})
		articles_dat[i]["id"] = article.ID
		articles_dat[i]["title"] = article.Title
		articles_dat[i]["author"] = article.Author
		articles_dat[i]["createdAt"] = article.CreatedAt.Format("02/01/2006 15:04:05")
	}
	main_app_dat["articles"] = articles_dat
	main_app_dat["more"] = len(articles) == 10
	main_app_dat["next"] = 2
	index_data["main_app"] = main_app_dat
	return c.Render(c.Response().Status, "index.html", index_data)
}

func GetFavicon(c echo.Context) error {
	return c.File("web/static/favicon.ico")
}

func GetBody(c echo.Context) error {
	var data Data
	data.Message = "Estoy seco pero el tequila ayuda"
	return c.Render(c.Response().Status, "body.html", data)
}

func GetNavbar(c echo.Context) error {
	var data Data
	userInfo := make(map[string]string)
	data.IsAuthenticated = false
	sess, _ := session.Get("session", c)
	if sess.Values["user_id"] != nil {
		data.IsAuthenticated = true
		user, err := database.GetUserByID(sess.Values["user_id"].(uint64))
		if err != nil {
			data.IsAuthenticated = false
			return c.Render(c.Response().Status, "navbar.html", data)
		}
		userInfo["username"] = user.Username
		if user.Firstname != nil && user.Lastname != nil {
			userInfo["fullname"] = *user.Firstname + " " + *user.Lastname
		}
	}
	data.UserInfo = userInfo
	return c.Render(c.Response().Status, "navbar", data)
}
