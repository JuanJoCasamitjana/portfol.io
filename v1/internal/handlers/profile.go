package handlers

import (
	"net/http"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/labstack/echo/v4"
)

func GetMyProfileInfo(c echo.Context) error {
	user, err := getUserOfSession(c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", "Error getting user")
	}
	profile, err := database.GetProfileByUserID(user.ID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", "Error getting profile")
	}
	data := map[string]interface{}{
		"username":  user.Username,
		"fullName":  user.FullName,
		"bio":       profile.Bio,
		"pfp":       profile.ImageURL,
		"pfp_thumb": profile.ThumbURL,
	}
	return c.Render(http.StatusOK, "profile", data)
}

func GetProfileEditForm(c echo.Context) error {
	user, err := getUserOfSession(c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", "Error getting user")
	}
	profile, err := database.GetProfileByUserID(user.ID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", "Error getting profile")
	}
	data := map[string]interface{}{
		"username": user.Username,
		"bio":      profile.Bio,
	}
	return c.Render(http.StatusOK, "profile_edit", data)
}

func EditProfileInfo(c echo.Context) error {
	user, err := getUserOfSession(c)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", "Error getting user")
	}
	bio := c.FormValue("bio")
	pfp, err := c.FormFile("pfp")
	data := make(map[string]interface{})
	data["bio"] = bio
	data["username"] = user.Username
	if err != nil {
		data["error"] = "Error getting pfp file"
		return c.Render(http.StatusOK, "profile_edit", data)
	}
	img, err := convertFileToBytes(pfp)
	if err != nil {
		data["error"] = "Error converting pfp file"
		return c.Render(http.StatusOK, "profile_edit", data)
	}
	urls, err := uploadImageToImgbb(img)
	if err != nil {
		data["error"] = "Error uploading pfp file"
		return c.Render(http.StatusOK, "profile_edit", data)
	}
	profile, err := database.GetProfileByUserID(user.ID)
	if err != nil {
		data["error"] = "Error getting profile"
		return c.Render(http.StatusOK, "profile_edit", data)
	}
	profile.Bio = bio
	profile.ImageURL = urls["image_url"]
	profile.ThumbURL = urls["thumb_url"]
	profile.DeleteURL = urls["delete_url"]
	err = database.UpdateProfile(&profile)
	if err != nil {
		data["error"] = "Error updating profile"
		return c.Render(http.StatusOK, "profile_edit", data)
	}
	data = map[string]interface{}{
		"username":  user.Username,
		"fullName":  user.FullName,
		"bio":       profile.Bio,
		"pfp":       profile.ImageURL,
		"pfp_thumb": profile.ThumbURL,
	}
	return c.Render(http.StatusOK, "profile", data)
}
