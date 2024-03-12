package handlers

import (
	"errors"
	"log"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/JuanJoCasamitjana/portfol.io/internal/utils"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetUserOfSession(c echo.Context) (model.User, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return model.User{}, err
	}
	userID, ok := sess.Values["user_id"].(uint64)
	if !ok {
		return model.User{}, errors.New("user not found in session")
	}
	user, err := database.FindUserById(userID)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// Registration process, linked only to register.html template and locale/*/register.json
func GetRegisterForm(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "register", data)
}

func Register(c echo.Context) error {
	locale := utils.GetLocale(c)
	var user model.User
	username, fullname := c.FormValue("username"), c.FormValue("fullname")
	password, password2 := c.FormValue("password"), c.FormValue("password2")
	form_errors := make(map[string]string)
	form_values := map[string]string{
		"username": username,
		"fullname": fullname,
	}
	if password != password2 {
		form_errors["password2"] = utils.Translate(locale, "register_password_mismatch_error")
	}
	err := user.ValidateAndSetUsername(username)
	_, err_db := database.FindUserByUsername(username)
	if err != nil {
		form_errors["username"] = utils.Translate(locale, "register_username_invalid_error") + ". "
	}
	if err_db == nil {
		form_errors["username"] = form_errors["username"] + utils.Translate(locale, "register_username_taken_error")
	}
	err = user.Password.ValidateAndSetPassword(password)
	if err != nil {
		form_errors["password"] = utils.Translate(locale, "register_password_invalid_error")
	}
	if len(form_errors) > 0 {
		data := map[string]any{
			"locale":     locale,
			"errors":     form_errors,
			"formValues": form_values,
		}
		return c.Render(200, "register", data)
	}
	user.FullName = fullname
	err = database.CreateUser(&user)
	if err != nil {
		form_errors["other"] = "register_user_creation_error"
		data := map[string]any{
			"locale":     locale,
			"errors":     form_errors,
			"formValues": form_values,
		}
		return c.Render(200, "register", data)
	}
	err = login_user_session(user, c)
	if err != nil {
		log.Println(err)
		return c.Render(200, "success", nil)
	}
	//Notify HTMX that the session has changed
	c.Response().Header().Set("HX-Trigger", "session-changed")
	return c.Render(200, "success", nil)
}

// Login process, linked only to login.html template and locale/*/login.json

func GetLoginForm(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "login", data)
}

func Login(c echo.Context) error {
	locale := utils.GetLocale(c)
	username, password := c.FormValue("username"), c.FormValue("password")
	form_errors := make(map[string]string)
	form_values := map[string]string{
		"username": username,
	}
	user, err := database.FindUserByUsername(username)
	if err != nil {
		form_errors["username"] = utils.Translate(locale, "login_username_not_found_error")
		data := map[string]any{
			"locale":     locale,
			"errors":     form_errors,
			"formValues": form_values,
		}
		return c.Render(200, "login", data)
	}
	if !user.Password.ComparePassword(password) {
		form_errors["password"] = utils.Translate(locale, "login_password_mismatch_error")
		data := map[string]any{
			"locale":     locale,
			"errors":     form_errors,
			"formValues": form_values,
		}
		return c.Render(200, "login", data)
	}
	err = login_user_session(user, c)
	if err != nil {
		form_errors["other"] = utils.Translate(locale, "login_session_error")
		data := map[string]any{
			"locale":     locale,
			"errors":     form_errors,
			"formValues": form_values,
		}
		return c.Render(200, "login", data)
	}
	//Notify HTMX that the session has changed
	c.Response().Header().Set("HX-Trigger", "session-changed")
	return c.Render(200, "success", nil)
}

func login_user_session(user model.User, c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Values["user_id"] = user.ID
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}
	return nil
}

// Logout is simple, just remove the user_id from the session
func Logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	delete(sess.Values, "user_id")
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return err
	}
	//Notify HTMX that the session has changed
	c.Response().Header().Set("HX-Trigger", "session-changed")
	return c.Render(200, "success", nil)
}
