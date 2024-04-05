package handlers

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

func IsAdmin(c echo.Context) bool {
	user, err := GetUserOfSession(c)
	if err != nil {
		return false
	}
	return user.Authority.Level >= model.AUTH_ADMIN.Level
}

func IsModerator(c echo.Context) bool {
	user, err := GetUserOfSession(c)
	if err != nil {
		return false
	}
	return user.Authority.Level >= model.AUTH_MODERATOR.Level
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
	user := model.NewUser()
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
		"isActive":        user.Active,
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
		"current_avatar":  user.Profile.PfPUrl,
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
	current_avatar := c.FormValue("current_avatar")
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
		avatar_url := current_avatar
		if urls["thumb_url"] != "" {
			avatar_url = urls["thumb_url"]
		}
		data := map[string]any{
			"locale":         locale,
			"errors":         form_errors,
			"username":       user.Username,
			"bio":            bio,
			"email":          email,
			"fullname":       fullname,
			"avatar":         urls["thumb_url"],
			"current_avatar": avatar_url,
		}
		return c.Render(200, "profile_edit", data)
	}
	user.Profile.PfPUrl = urls["thumb_url"]
	if urls["thumb_url"] == "" {
		user.Profile.PfPUrl = current_avatar
	}
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
		"isActive":        user.Active,
	}
	return c.Render(200, "profile", data)
}

func isValidEmail(email string, current_user model.User) bool {
	//False if the email is already in use
	if email == "" {
		return true
	}
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
	session_user, _ := GetUserOfSession(c)
	is_current_user := session_user.Username == user.Username
	user_follow_list, err := database.FindFollowListByUsername(session_user.Username)
	if err != nil {
		user_follow_list = model.FollowList{Owner: session_user.Username}
		err = database.CreateFollowList(&user_follow_list)
		if err != nil {
			return c.String(500, "Internal server error")
		}
	}
	is_following := isFollowing(user_follow_list, user)
	data := map[string]any{
		"username":        user.Username,
		"fullname":        user.FullName,
		"locale":          utils.GetLocale(c),
		"bio":             user.Profile.Bio,
		"avatar":          user.Profile.PfPUrl,
		"is_current_user": is_current_user,
		"is_following":    is_following,
		"isActive":        user.Active,
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
	if err != nil || !user.Active {
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
	if err != nil || !user.Active {
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
	if err != nil || !user.Active {
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
	if err != nil || !user.Active {
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
		"isActive": user.Active,
	}
	return c.Render(200, "sections_list", data)
}

func GetSectionEdit(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil || !user.Active {
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

func FollowUser(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	username := c.Param("username")
	if user.Username == username {
		return c.String(400, "Bad request")
	}
	user_to_be_followed, err := database.FindUserByUsername(username)
	if err != nil {
		return c.String(404, "Not found")
	}
	user_follow_list, err := database.FindFollowListByUsername(user.Username)
	if err != nil {
		user_follow_list = model.FollowList{Owner: user.Username}
		err = database.CreateFollowList(&user_follow_list)
		if err != nil {
			return c.String(500, "Internal server error")
		}
	}
	err = database.FollowUser(&user_follow_list, &user_to_be_followed)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	data := map[string]any{
		"username":     user_to_be_followed.Username,
		"locale":       locale,
		"is_following": true,
	}
	return c.Render(200, "follow_button", data)
}

func UnfollowUser(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	username := c.Param("username")
	if user.Username == username {
		return c.String(400, "Bad request")
	}
	user_to_be_unfollowed, err := database.FindUserByUsername(username)
	if err != nil {
		return c.String(404, "Not found")
	}
	user_follow_list, err := database.FindFollowListByUsername(user.Username)
	if err != nil {
		return c.String(404, "Not found")
	}
	err = database.UnfollowUser(&user_follow_list, &user_to_be_unfollowed)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	data := map[string]any{
		"username":     user_to_be_unfollowed.Username,
		"locale":       locale,
		"is_following": false,
	}
	return c.Render(200, "follow_button", data)
}

func isFollowing(follower_follow_list model.FollowList, followed model.User) bool {
	for _, user := range follower_follow_list.Following {
		if user.Username == followed.Username {
			return true
		}
	}
	return false
}

func FollowingPostsPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	postsDB, err := database.FindFollowingPostsPaginated(user, page, 12)
	if err != nil {
		return c.String(404, "Not found")
	}
	posts := convertPostsToDataMap(postsDB)
	more := len(posts) == 12
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/following?page=%d", page+1)
	}
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
		"posts":    posts,
		"more":     more,
		"nextPage": nextPageLoader,
	}
	return c.Render(200, "posts", data)
}

func ListWhoIFollow(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil {
		return c.String(401, "Unauthorized")
	}
	user_follow_list, err := database.FindFollowListByUsername(user.Username)
	if err != nil {
		user_follow_list = model.FollowList{Owner: user.Username}
		err = database.CreateFollowList(&user_follow_list)
		if err != nil {
			return c.String(500, "Internal server error")
		}
	}
	users := make([]map[string]any, len(user_follow_list.Following))
	for i, user := range user_follow_list.Following {
		users[i] = map[string]any{
			"username": user.Username,
			"fullname": user.FullName,
		}
	}
	data := map[string]any{
		"locale":    locale,
		"username":  user.Username,
		"following": users,
	}
	return c.Render(200, "following", data)
}

func GetModDashBoard(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_MODERATOR.Level || !user.Active {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
		"isAdmin":  user.Authority.Level == model.AUTH_ADMIN.Level,
	}
	return c.Render(200, "dashboard", data)
}

func GetAdminDashBoard(c echo.Context) error {
	user, err := GetUserOfSession(c)
	//An admin cannot be banned because it is the highest authority
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale":   locale,
		"username": user.Username,
		"isAdmin":  user.Authority.Level == model.AUTH_ADMIN.Level,
	}
	return c.Render(200, "dashboard", data)
}

func GetUsersListPaginated(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_MODERATOR.Level {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	usersDB, err := database.FindUsersPaginated(page, 12)
	if err != nil {
		return c.String(404, "Not found")
	}
	users := convertUsersToDataMap(usersDB, user.Authority.Level, locale)
	more := len(users) == 12
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/moderation/users?page=%d", page+1)
	}
	data := map[string]any{
		"locale":   locale,
		"users":    users,
		"more":     more,
		"nextPage": nextPageLoader,
	}
	if page == 1 {
		return c.Render(200, "users_list", data)
	}
	return c.Render(200, "users_list_page", data)

}

func GetUsersListSearchPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_MODERATOR.Level {
		return c.String(401, "Unauthorized")
	}
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	search := c.QueryParam("search")
	usersDB, err := database.FindUsersPaginatedBySearch(search, page, 12)
	if err != nil {
		return c.String(404, "Not found")
	}
	users := convertUsersToDataMap(usersDB, user.Authority.Level, locale)
	more := len(users) == 12
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/moderation/users/search?search=%s&page=%d", search, page+1)
	}
	data := map[string]any{
		"locale":   locale,
		"users":    users,
		"more":     more,
		"nextPage": nextPageLoader,
	}
	return c.Render(200, "users_list_page", data)
}

func convertUsersToDataMap(users []model.User, requester_level uint8, locale string) []map[string]any {
	users_list := make([]map[string]any, len(users))
	for i, user := range users {
		isActionaAvailable := requester_level > user.Authority.Level
		users_list[i] = map[string]any{
			"username": user.Username,
			"fullname": user.FullName,
			"email":    user.Email,
			"active":   user.Active,
			"auth":     user.Authority.AuthName,
			"avatar":   user.Profile.PfPUrl,
			"bio":      user.Profile.Bio,
			"locale":   locale,
			"action":   isActionaAvailable,
		}
	}
	return users_list
}

func BanUser(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_MODERATOR.Level || !user.Active {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	username := c.Param("username")
	user_to_be_banned, err := database.FindUserByUsername(username)
	if err != nil {
		return c.String(404, "Not found")
	}
	user_to_be_banned.Active = false
	err = database.UpdateUser(&user_to_be_banned)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	data := map[string]any{
		"username": user_to_be_banned.Username,
		"fullname": user_to_be_banned.FullName,
		"email":    user_to_be_banned.Email,
		"active":   user_to_be_banned.Active,
		"auth":     user_to_be_banned.Authority.AuthName,
		"avatar":   user_to_be_banned.Profile.PfPUrl,
		"bio":      user_to_be_banned.Profile.Bio,
		"locale":   locale,
		"action":   true,
	}
	return c.Render(200, "user_item", data)
}

func UnbanUser(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_MODERATOR.Level || !user.Active {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	username := c.Param("username")
	user_to_be_banned, err := database.FindUserByUsername(username)
	if err != nil {
		return c.String(404, "Not found")
	}
	user_to_be_banned.Active = true
	err = database.UpdateUser(&user_to_be_banned)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	data := map[string]any{
		"username": user_to_be_banned.Username,
		"fullname": user_to_be_banned.FullName,
		"email":    user_to_be_banned.Email,
		"active":   user_to_be_banned.Active,
		"auth":     user_to_be_banned.Authority.AuthName,
		"avatar":   user_to_be_banned.Profile.PfPUrl,
		"bio":      user_to_be_banned.Profile.Bio,
		"locale":   locale,
		"action":   true,
	}
	return c.Render(200, "user_item", data)
}

func GetConfigChangeForm(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale":                     locale,
		"imgbb_api_key":              IMGBB_API_KEY,
		"corporative_email":          utils.FromEmail,
		"corporative_email_password": string(utils.FromEmailPassword),
		"smtp_server":                utils.SmtpHost,
		"smtp_port":                  utils.SmtpPort,
	}
	return c.Render(200, "config_change", data)
}

func ChangeConfig(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	imgbb_api_key := c.FormValue("imgbb_api_key")
	corporative_email := c.FormValue("corporative_email")
	corporative_email_password := c.FormValue("corporative_email_password")
	smtp_server := c.FormValue("smtp_server")
	smtp_port := c.FormValue("smtp_port")
	form_errors := make(map[string]string)
	if imgbb_api_key == "" {
		form_errors["imgbb_api_key"] = utils.Translate(locale, "config_change_imgbb_api_key_error")
	}
	if corporative_email == "" {
		form_errors["corporative_email"] = utils.Translate(locale, "config_change_corporative_email_error")
	}
	if !isValidEmail(corporative_email, user) {
		form_errors["corporative_email"] = utils.Translate(locale, "config_change_corporative_email_invalid_error")
	}
	if corporative_email_password == "" {
		form_errors["corporative_email_password"] = utils.Translate(locale, "config_change_corporative_email_password_error")
	}
	if smtp_server == "" {
		form_errors["smtp_server"] = utils.Translate(locale, "config_change_smtp_server_error")
	}
	if smtp_port == "" {
		form_errors["smtp_port"] = utils.Translate(locale, "config_change_smtp_port_error")
	}
	if len(form_errors) > 0 {
		data := map[string]any{
			"locale":                     locale,
			"errors":                     form_errors,
			"imgbb_api_key":              imgbb_api_key,
			"corporative_email":          corporative_email,
			"corporative_email_password": corporative_email_password,
			"smtp_server":                smtp_server,
			"smtp_port":                  smtp_port,
		}
		return c.Render(200, "config_change", data)
	}
	IMGBB_API_KEY = imgbb_api_key
	utils.FromEmail = corporative_email
	utils.FromEmailPassword = []byte(corporative_email_password)
	utils.SmtpHost = smtp_server
	utils.SmtpPort = smtp_port
	data := map[string]any{
		"locale":  locale,
		"message": utils.Translate(locale, "config_change_success"),
	}
	return c.Render(200, "config_change", data)
}

func CreateNewModeratorForm(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "moderator_new", data)
}

func CreateNewModerator(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	username := c.FormValue("username")
	password := c.FormValue("password")
	password2 := c.FormValue("password_2")
	form_errors := make(map[string]string)
	_, err = database.FindUserByUsername(username)
	if err == nil {
		form_errors["username"] = utils.Translate(locale, "moderator_new_username_taken_error")
	}
	new_mod := model.NewUser()
	err = new_mod.ValidateAndSetUsername(username)
	if err != nil {
		form_errors["username"] = utils.Translate(locale, "moderator_new_username_invalid_error")
	}
	err = new_mod.Password.ValidateAndSetPassword(password)
	if err != nil {
		form_errors["password"] = utils.Translate(locale, "moderator_new_password_invalid_error")
	}
	if password != password2 {
		form_errors["password_2"] = utils.Translate(locale, "moderator_new_password_mismatch_error")
	}
	if len(form_errors) > 0 {
		data := map[string]any{
			"locale":   locale,
			"errors":   form_errors,
			"username": username,
		}
		return c.Render(200, "moderator_new", data)
	}
	new_mod.Authority = model.AUTH_MODERATOR
	err = database.CreateUser(&new_mod)
	if err != nil {
		form_errors["other"] = utils.Translate(locale, "moderator_new_error")
		data := map[string]any{
			"locale":   locale,
			"errors":   form_errors,
			"username": username,
		}
		return c.Render(200, "moderator_new", data)
	}
	data := map[string]any{
		"locale":  locale,
		"message": utils.Translate(locale, "moderator_new_success"),
	}
	return c.Render(200, "moderator_new", data)
}

func ShowApplicationSummary(c echo.Context) error {
	locale := utils.GetLocale(c)
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_MODERATOR.Level {
		return c.String(401, "Unauthorized")
	}
	userCount, err := database.CountUsers()
	if err != nil {
		return c.String(500, "Internal server error")
	}
	galleryCount, err := database.CountGalleries()
	if err != nil {
		return c.String(500, "Internal server error")
	}
	articleCount, err := database.CountArticles()
	if err != nil {
		return c.String(500, "Internal server error")
	}

	data := map[string]any{
		"locale":         locale,
		"articleCount":   articleCount,
		"galleryCount":   galleryCount,
		"userCount":      userCount,
		"totalPostCount": articleCount + galleryCount,
	}
	return c.Render(200, "application_summary", data)
}

func GetUserSearch(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "user_search", data)
}

func UserSearchPaginated(c echo.Context) error {
	locale := utils.GetLocale(c)
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	search := c.QueryParam("query")
	usersDB, err := database.FindUsersPaginatedBySearch(search, page, 12)
	if err != nil {
		return c.String(404, "Not found")
	}
	users := convertUsersToDataMap(usersDB, 0, locale)
	more := len(users) == 12
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/users/search?query=%s&page=%d", search, page+1)
	}
	data := map[string]any{
		"locale":   locale,
		"users":    users,
		"more":     more,
		"nextPage": nextPageLoader,
	}
	return c.Render(200, "user_list", data)
}

func GetRestraintAccessForm(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale":             locale,
		"isAccessRestricted": IsAccessRestricted,
	}
	return c.Render(200, "restraint_access", data)
}

func RestrainAccess(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	locale := utils.GetLocale(c)
	IsAccessRestricted = !IsAccessRestricted
	data := map[string]any{
		"locale":             locale,
		"isAccessRestricted": IsAccessRestricted,
	}
	return c.Render(200, "restraint_access", data)
}

func SendCopyOfDB(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	return c.File("./" + database.DBname)
}

func SendCopyOfLogs(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_ADMIN.Level {
		return c.String(401, "Unauthorized")
	}
	momment := strconv.FormatInt(time.Now().Unix(), 10)
	name := "logs_" + momment + ".zip"
	err = CompressLogs(name)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	return c.File("./" + name)
}

func CompressLogs(name string) error {
	zipfile, err := os.Create(name)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	zipWriter := zip.NewWriter(zipfile)
	defer zipWriter.Close()
	logs_folder := "./logs"

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		f, err := zipWriter.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
		return nil
	}
	err = filepath.Walk(logs_folder, walker)
	if err != nil {
		return err
	}
	return nil
}
