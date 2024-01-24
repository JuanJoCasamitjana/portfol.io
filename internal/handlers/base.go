package handlers

import (
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
	ErrorsExist     bool
	MessageExist    bool
}

func GetIndex(c echo.Context) error {
	var data Data
	data.Title = "Portfol.io"
	return c.Render(c.Response().Status, "index.html", data)
}

func GetFavicon(c echo.Context) error {
	return c.File("web/static/favicon.ico")
}

func GetBody(c echo.Context) error {
	var data Data
	data.Message = "Welcome"
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
	return c.Render(c.Response().Status, "navbar.html", data)
}
