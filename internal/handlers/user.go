package handlers

import (
	"net/http"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetUserByUsername(c echo.Context) error {
	username := c.Param("username")
	user, err := database.GetUserByUsername(username)
	if err != nil {
		return c.Render(http.StatusNotFound, "404", nil)
	}
	return c.Render(http.StatusOK, "user", user)
}

func RegisterUser(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

func CreateUser(c echo.Context) error {
	var data Data
	var user model.User
	errorsmap := make(map[string]string)
	formValuesmap := make(map[string]string)
	username := c.FormValue("username")
	password := c.FormValue("password")
	firstName := c.FormValue("firstName")
	lastName := c.FormValue("lastName")
	passwordConfirmation := c.FormValue("password2")
	_, err := database.GetUserByUsername(username)
	formValuesmap["firstName"] = firstName
	formValuesmap["lastName"] = lastName
	if err == nil {
		errorsmap["username"] = "Username already taken. "
	} else {
		formValuesmap["username"] = username
		data.FormValues = formValuesmap
	}
	if username == "" {
		errorsmap["username"] = "Username cannot be empty. "
	}
	if err := user.ValidateAndSetUsername(username); err == model.ErrInvalidUsername {
		errorsmap["username"] = errorsmap["username"] + "Username must only contain characters from A-Z, a-z, 0-9 and [-,_,.]. "
	}
	if password != passwordConfirmation {
		errorsmap["password2"] = "Passwords do not match. "
	}
	if len(password) < 12 {
		errorsmap["password"] = "Password must be at least 12 characters long and at most 72 characters long. "
	}
	if len(password) > 72 {
		errorsmap["password"] = "Password must be at least 12 characters long and at most 72 characters long. "
	}
	err = user.Password.ValidateAndSetPassword(password)
	if err == model.ErrPasswordTooLong {
		errorsmap["password"] = "Password must be at least 12 characters long and at most 72 characters long. "
	}
	if err == model.ErrPasswordContainsUnsuportedCharacters {
		errorsmap["password"] = errorsmap["password"] + "Password must only contain caharcters from A-Z a-z 0-9 !#$%&()*+,-.:;<=>?@[]_{} and spaces. "
	}
	//Fallback error message
	if err != nil && err != model.ErrPasswordTooLong && err != model.ErrPasswordContainsUnsuportedCharacters {
		errorsmap["password"] = "Something went worng while hashing your password, please try again later. "
	}
	if len(firstName) > 256 {
		errorsmap["firstName"] = "First name is too long. "
	}
	if len(lastName) > 256 {
		errorsmap["lastName"] = "Last name is too long. "
	}
	if len(errorsmap) > 0 {
		data.Errors = errorsmap
		data.ErrorsExist = true
		return c.Render(http.StatusOK, "register.html", data)
	}
	if firstName != "" {
		user.Firstname = &firstName
	}
	if lastName != "" {
		user.Lastname = &lastName
	}
	err = database.SaveUser(&user)
	if err != nil {
		errorsmap["other"] = "Something went wrong while saving your user, please try again later. "
	}
	if len(errorsmap) > 0 {
		data.Errors = errorsmap
		data.ErrorsExist = true
		return c.Render(http.StatusOK, "register.html", data)
	}
	data.Message = "User created successfully."
	return c.Render(http.StatusOK, "success.html", data)
}

func GetLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func Login(c echo.Context) error {
	var data Data
	username := c.FormValue("username")
	password := c.FormValue("password")
	errorsmap := make(map[string]string)
	formValuesmap := make(map[string]string)
	formValuesmap["username"] = username
	data.FormValues = formValuesmap
	user, err := database.GetUserByUsername(username)
	if err != nil {
		errorsmap["username"] = "Username does not exist. "
		data.Errors = errorsmap
		data.ErrorsExist = true
		return c.Render(http.StatusOK, "login.html", data)
	}
	ok := user.Password.ComparePassword(password)
	if !ok {
		errorsmap["password"] = "Password is incorrect. "
		data.Errors = errorsmap
		data.ErrorsExist = true
		return c.Render(http.StatusOK, "login.html", data)
	}
	sess, err := session.Get("session", c)
	if err != nil {
		errorsmap["other"] = "Something went wrong while creating your session, please try again later. "
		data.Errors = errorsmap
		return c.Render(http.StatusOK, "login.html", data)
	}
	sess.Values["user_id"] = user.ID
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		errorsmap["other"] = "Something went wrong while creating your session, please try again later. "
		data.Errors = errorsmap
		return c.Render(http.StatusOK, "login.html", data)
	}
	data.Message = "Login successful."
	c.Response().Header().Set("HX-Trigger", "session-changed")
	return c.Render(http.StatusOK, "success.html", data)
}

func Logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusOK, "logout.html", nil)
	}
	sess.Values["user_id"] = nil
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.Render(http.StatusOK, "logout.html", nil)
	}
	c.Response().Header().Set("HX-Trigger", "session-changed")
	return c.Render(http.StatusOK, "logout.html", nil)
}
