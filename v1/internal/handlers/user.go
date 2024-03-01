package handlers

import (
	"log"
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
	data := map[string]string{
		"username": user.Username,
		"fullname": user.FullName,
	}
	return c.Render(http.StatusOK, "user_profile", data)
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
	fullname := c.FormValue("fullname")
	passwordConfirmation := c.FormValue("password2")
	_, err := database.GetUserByUsername(username)
	formValuesmap["fullname"] = fullname
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
	if len(fullname) > 512 {
		errorsmap["firstName"] = "Full name is too long. "
	}
	if len(errorsmap) > 0 {
		data.Errors = errorsmap
		data.ErrorsExist = true
		return c.Render(http.StatusOK, "register.html", data)
	}
	if fullname != "" {
		user.FullName = fullname
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
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusOK, "login.html", data)
	}
	sess.Values["user_id"] = user.ID
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.Render(http.StatusOK, "login.html", data)
	}
	c.Response().Header().Set("HX-Trigger", "session-changed")
	return c.Render(http.StatusOK, "success", data)
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
	return c.Render(http.StatusOK, "success", data)
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
	return c.Render(http.StatusOK, "success", nil)
}

func GetEditUser(c echo.Context) error {
	sess, err := session.Get("session", c)
	userID := sess.Values["user_id"].(uint64)
	if err != nil {
		return c.Render(http.StatusUnauthorized, "401.html", nil)
	}
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
	data := Data{
		FormValues: map[string]string{
			"username": user.Username,
			"fullname": user.FullName,
		},
	}
	return c.Render(http.StatusOK, "edit_user", data)
}

func GetMyProfile(c echo.Context) error {
	sess, err := session.Get("session", c)
	userID := sess.Values["user_id"].(uint64)
	if err != nil {
		return c.Render(http.StatusUnauthorized, "401.html", nil)
	}
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
	data := map[string]string{
		"username": user.Username,
		"fullname": user.FullName,
	}
	return c.Render(http.StatusOK, "user", data)
}

func GetPasswordEdit(c echo.Context) error {
	return c.Render(http.StatusOK, "change_password", nil)
}

func ChangePassword(c echo.Context) error {
	var data Data
	sess, err := session.Get("session", c)
	userID := sess.Values["user_id"].(uint64)
	if err != nil {
		return c.Render(http.StatusUnauthorized, "401.html", nil)
	}
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusNotFound, "404.html", nil)
	}
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("password")
	newPasswordConfirmation := c.FormValue("password2")
	errorsmap := make(map[string]string)
	ok := user.Password.ComparePassword(oldPassword)
	if !ok {
		errorsmap["old_password"] = "Password is incorrect. "
	}
	if newPassword != newPasswordConfirmation {
		errorsmap["password2"] = "Passwords do not match. "
	}
	if len(newPassword) < 12 {
		errorsmap["password"] = "Password must be at least 12 characters long and at most 72 characters long. "
	}
	if len(newPassword) > 72 {
		errorsmap["password"] = "Password must be at least 12 characters long and at most 72 characters long. "
	}
	err = user.Password.ValidateAndSetPassword(newPassword)
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
	if len(errorsmap) > 0 {
		data.Errors = errorsmap
		data.ErrorsExist = true
		return c.Render(http.StatusOK, "change_password", data)
	}
	err = database.SaveUser(user)
	if err != nil {
		errorsmap["other"] = "Something went wrong while saving your user, please try again later. "
		data.Errors = errorsmap
		data.ErrorsExist = true
		return c.Render(http.StatusOK, "change_password", data)
	}
	data.Message = "Password changed successfully."
	return c.Render(http.StatusOK, "success", data)
}

func DeleteUser(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	userID := sess.Values["user_id"].(uint64)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	articles, _ := database.GetAllArticlesOfUser(user.Username)
	imgColls, _ := database.GetAllImageCollectionsOfUserWithImages(user.Username)
	for _, article := range articles {
		_ = database.DeleteArticle(&article)
	}
	for _, imgColl := range imgColls {
		for _, img := range imgColl.Images {
			_ = database.DeleteImage(&img)
		}
		_ = database.DeleteImageCollection(&imgColl)
	}
	log.Println("Delete user: ", user.Username, " id: ", user.ID)
	err = database.DeleteUser(user)
	if err != nil {
		log.Println("Error deleting user: ", err, " id: ", user.ID)
		return c.Render(http.StatusInternalServerError, "error", nil)
	}
	sess.Values["user_id"] = nil
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Panicln("Error saving session: ", err)
	}
	c.Response().Header().Set("HX-Trigger", "session-changed")
	return c.Render(http.StatusOK, "success", nil)
}
