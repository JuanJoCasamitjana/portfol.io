package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"strconv"

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

func ChangePasswordForm(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "password_change", data)
}

func ChangePassword(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	old_password, new_password, new_password2 := c.FormValue("old_password"), c.FormValue("new_password"), c.FormValue("new_password_2")
	form_errors := make(map[string]string)
	if !user.Password.ComparePassword(old_password) {
		form_errors["old_password"] = utils.Translate(locale, "password_change_old_password_error")
	}
	if new_password != new_password2 {
		form_errors["new_password_2"] = utils.Translate(locale, "password_change_new_password_mismatch_error")
	}
	err = user.Password.ValidateAndSetPassword(new_password)
	if err != nil {
		form_errors["new_password"] = utils.Translate(locale, "password_change_new_password_invalid_error")
	}
	if len(form_errors) > 0 {
		data := map[string]any{
			"locale": locale,
			"errors": form_errors,
		}
		return c.Render(200, "password_change", data)
	}
	err = database.UpdateUser(&user)
	if err != nil {
		form_errors["other"] = utils.Translate(locale, "password_change_error")
		data := map[string]any{
			"locale": locale,
			"errors": form_errors,
		}
		return c.Render(200, "password_change", data)
	}
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "success", data)
}

func DeleteProfile(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.Render(401, "error", nil)
	}
	err = database.DeleteUser(&user)
	if err != nil {
		return c.Render(500, "error", nil)
	}
	c.Response().Header().Set("HX-Trigger", "session-changed")
	return c.Render(200, "success", nil)
}

func GetMyProfile(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.Render(401, "error", nil)
	}
	data := map[string]any{
		"username":        user.Username,
		"fullname":        user.FullName,
		"locale":          utils.GetLocale(c),
		"bio":             user.Profile.Bio,
		"avatar":          user.Profile.PfPUrl,
		"email":           user.Email,
		"is_current_user": true,
	}
	return c.Render(200, "profile", data)
}

func GetProfileEditForm(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.Render(401, "error", nil)
	}
	data := map[string]any{
		"username":        user.Username,
		"fullname":        user.FullName,
		"locale":          utils.GetLocale(c),
		"bio":             user.Profile.Bio,
		"avatar":          user.Profile.PfPUrl,
		"email":           user.Email,
		"is_current_user": true,
	}
	return c.Render(200, "profile_edit", data)
}

func EditProfile(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.Render(401, "error", nil)
	}
	locale := utils.GetLocale(c)
	bio, email, fullname := c.FormValue("bio"), c.FormValue("email"), c.FormValue("fullname")
	var form_errors = map[string]string{}
	if !isValidEmail(email, user) {
		form_errors["email"] = utils.Translate(locale, "profile_edit_email_invalid_error")
	}
	avatar, err := c.FormFile("avatar")
	var urls = map[string]string{"thumb_url": "", "delete_url": ""}
	if err == nil && avatar.Size > 0 {
		avatar_bytes, err := convertFileToBytes(avatar)
		if err != nil {
			form_errors["avatar"] = utils.Translate(locale, "profile_edit_avatar_server_error")
		}
		urls, err = uploadImageToImgbb(avatar_bytes)
		if err != nil {
			form_errors["avatar"] = utils.Translate(locale, "profile_edit_avatar_server_error")
		}
	}
	if len(form_errors) > 0 {
		data := map[string]any{
			"locale":   locale,
			"errors":   form_errors,
			"username": user.Username,
			"bio":      bio,
			"email":    email,
			"fullname": fullname,
			"avatar":   urls["thumb_url"],
		}
		return c.Render(200, "profile_edit", data)
	}
	user.Profile.PfPUrl = urls["thumb_url"]
	user.Profile.PfPDeleteUrl = urls["delete_url"]
	user.Profile.Bio = bio
	user.Email = email
	user.FullName = fullname
	err = database.UpdateUser(&user)
	if err != nil {
		form_errors["other"] = "profile_edit_error"
		data := map[string]any{
			"locale":   locale,
			"username": user.Username,
			"errors":   form_errors,
			"bio":      bio,
			"email":    email,
			"fullname": fullname,
			"avatar":   urls["thumb_url"],
		}
		return c.Render(200, "profile_edit", data)
	}
	data := map[string]any{
		"locale":          locale,
		"username":        user.Username,
		"fullname":        user.FullName,
		"bio":             user.Profile.Bio,
		"avatar":          user.Profile.PfPUrl,
		"email":           user.Email,
		"is_current_user": true,
	}
	return c.Render(200, "profile", data)
}

func isValidEmail(email string, current_user model.User) bool {
	//False if the email is already in use
	user, err := database.FindUserByEmail(email)
	if err != nil {
		return true
	}
	if user.ID != current_user.ID {
		return false
	}
	_, err = mail.ParseAddress(email)
	return err == nil
}

func GetUserProfile(c echo.Context) error {
	username := c.Param("username")
	user, err := database.FindUserByUsername(username)
	if err != nil {
		return c.Render(404, "error", nil)
	}
	data := map[string]any{
		"username":        user.Username,
		"fullname":        user.FullName,
		"locale":          utils.GetLocale(c),
		"bio":             user.Profile.Bio,
		"avatar":          user.Profile.PfPUrl,
		"is_current_user": false,
	}
	return c.Render(200, "profile", data)
}

func GetUserSections(c echo.Context) error {
	username := c.Param("username")
	locale := utils.GetLocale(c)
	mainSection := username
	sections, err := database.FindSectionsByUser(username)
	if err != nil {
		return c.Render(404, "error", nil)
	}
	sections_list := []string{mainSection}
	for _, section := range sections {
		sections_list = append(sections_list, section.Name)
	}
	data := map[string]any{
		"locale":   locale,
		"username": username,
		"section":  mainSection,
		"sections": sections_list,
	}
	return c.Render(200, "sections", data)
}

func GetUserSectionPaginated(c echo.Context) error {
	username, section_name := c.Param("username"), c.Param("section")
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	locale := utils.GetLocale(c)
	mainSection := username
	if section_name == mainSection {
		postsDB, err := database.FindPostsByUserPaginated(username, page, 12)
		if err != nil {
			return c.String(404, "Not found")
		}
		posts := convertPostsToDataMap(postsDB)
		more := len(posts) == 12
		nextPageLoader := ""
		if more {
			nextPageLoader = fmt.Sprintf("/profile/%s/section/%s?page=%d", username, section_name, page+1)
		}
		data := map[string]any{
			"locale":   locale,
			"username": username,
			"posts":    posts,
			"more":     more,
			"nextPage": nextPageLoader,
		}

		return c.Render(200, "posts", data)
	}
	postsDB, err := database.FindPostsByUserAndSectionPaginated(username, section_name, page, 12)
	if err != nil {
		return c.String(404, "Not found")
	}
	posts := convertPostsToDataMap(postsDB)
	more := len(posts) == 12
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/profile/%s/section/%s?page=%d", username, section_name, page+1)
	}
	data := map[string]any{
		"locale":   locale,
		"username": username,
		"posts":    posts,
		"more":     more,
		"nextPage": nextPageLoader,
	}
	return c.Render(200, "posts", data)
}

func CreateNewSectionForm(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
	}
	return c.Render(200, "section_new", data)
}

func CreateNewSection(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	section_name := c.FormValue("name")
	form_errors := make(map[string]string)
	if len(section_name) < 5 {
		form_errors["name"] = utils.Translate(locale, "section_new_name_short_error") + ". "
	}
	if section_name == user.Username {
		form_errors["name"] = form_errors["name"] + utils.Translate(locale, "section_new_name_invalid_error") + ". "
	}
	section, err := database.FindSectionByUsernameAndName(user.Username, section_name)
	if err == nil {
		form_errors["name"] = form_errors["name"] + utils.Translate(locale, "section_new_name_taken_error") + ". "
	}
	if len(form_errors) > 0 {
		data := map[string]any{
			"locale":   locale,
			"username": user.Username,
			"name":     section_name,
			"errors":   form_errors,
		}
		return c.Render(200, "section_new", data)
	}
	section = model.Section{Name: section_name, Owner: user.Username}
	err = database.CreateSection(&section)
	if err != nil {
		data := map[string]any{
			"locale":   locale,
			"username": user.Username,
			"name":     section_name,
			"errors":   map[string]string{"other": utils.Translate(locale, "section_new_error")},
		}
		return c.Render(200, "section_new", data)
	}
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
		"name":     section_name,
	}
	return c.Render(200, "section_edit", data)
}

func AddPostToSection(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	section_name := c.Param("section")
	section, err := database.FindSectionByUsernameAndName(user.Username, section_name)
	if err != nil {
		return c.String(404, "Not found")
	}
	postIDstr := c.Param("post")
	postID, err := strconv.ParseUint(postIDstr, 10, 64)
	if err != nil {
		return c.String(404, "Not found")
	}
	post, err := database.FindPostById(postID)
	if err != nil {
		return c.String(404, "Not found")
	}
	err = database.AddPostToSection(&section, &post)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	c.Response().Header().Set("HX-Trigger", "section-changed")
	return c.String(200, "OK")
}

func RemovePostFromSection(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	section_name := c.Param("section")
	section, err := database.FindSectionByUsernameAndName(user.Username, section_name)
	if err != nil {
		return c.String(404, "Not found")
	}
	postIDstr := c.Param("post")
	postID, err := strconv.ParseUint(postIDstr, 10, 64)
	if err != nil {
		return c.String(404, "Not found")
	}
	post, err := database.FindPostById(postID)
	if err != nil {
		return c.String(404, "Not found")
	}
	err = database.RemovePostFromSection(&section, &post)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	c.Response().Header().Set("HX-Trigger", "section-changed")
	return c.String(200, "OK")
}

func GetModificablePostsFromSectionPaginated(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	section_name := c.Param("section")
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	locale := utils.GetLocale(c)
	postsDB, err := database.FindPostsByUserAndSectionPaginated(user.Username, section_name, page, 12)
	if err != nil {
		return c.String(404, "Not found")
	}
	posts := convertPostsToDataMap(postsDB)
	more := len(posts) == 12
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/profile/%s/section/%s?page=%d", user.Username, section_name, page+1)
	}
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
		"section":  section_name,
		"posts":    posts,
		"more":     more,
		"nextPage": nextPageLoader,
		"add":      false,
	}
	return c.Render(200, "posts_modificable", data)
}

func GetModificablePostsNotFromSectionPaginated(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	section_name := c.Param("section")
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	locale := utils.GetLocale(c)
	postsDB, err := database.FindPostsByUserNotInSectionPaginated(user.Username, section_name, page, 12)
	if err != nil {
		return c.String(404, "Not found")
	}
	posts := convertPostsToDataMap(postsDB)
	more := len(posts) == 12
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/profile/%s/section/%s?page=%d", user.Username, section_name, page+1)
	}
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
		"section":  section_name,
		"posts":    posts,
		"more":     more,
		"nextPage": nextPageLoader,
		"add":      true,
	}
	return c.Render(200, "posts_modificable", data)
}

func DeleteSection(c echo.Context) error {
	section, username := c.Param("section"), c.Param("username")
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	if user.Username != username {
		return c.String(401, "Unauthorized")
	}
	err = database.DeleteSectionByUsernameAndName(username, section)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	sections, err := database.FindSectionsByUser(username)
	if err != nil {
		return c.String(404, "Not found")
	}
	sections_lists := make([]string, len(sections))
	for i, section := range sections {
		sections_lists[i] = section.Name
	}
	data := map[string]any{
		"locale":   utils.GetLocale(c),
		"username": username,
		"sections": sections_lists,
	}
	return c.Render(200, "sections_list", data)
}

func GetMySectionsList(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	sections, err := database.FindSectionsByUser(user.Username)
	if err != nil {
		return c.String(404, "Not found")
	}
	sections_lists := make([]string, len(sections))
	for i, section := range sections {
		sections_lists[i] = section.Name
	}
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
		"sections": sections_lists,
	}
	return c.Render(200, "sections_list", data)
}

func GetSectionEdit(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	section_name := c.Param("section")
	section, err := database.FindSectionByUsernameAndName(user.Username, section_name)
	if err != nil {
		return c.String(404, "Not found")
	}
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
		"name":     section.Name,
	}
	return c.Render(200, "section_edit", data)
}
