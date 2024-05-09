package handlers

import (
	"fmt"
	"strconv"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/JuanJoCasamitjana/portfol.io/internal/utils"
	"github.com/labstack/echo/v4"
)

func GetCreateReport(c echo.Context) error {
	which := c.QueryParam("which")
	if which == "part" {
		return GetCreateReportPart(c)
	}
	return GetCreateReportFull(c)
}

func GetCreateReportFull(c echo.Context) error {
	user, err := GetUserOfSession(c)
	locale := utils.GetLocale(c)
	isAuthenticated := err == nil
	isModerator := IsModerator(c)
	isAdmin := IsAdmin(c)
	data := map[string]any{
		"app_title":       "Portfol.io",
		"locale":          locale,
		"isActive":        user.Active,
		"IsAuthenticated": isAuthenticated,
		"IsModerator":     isModerator,
		"IsAdmin":         isAdmin,
		"page_to_load":    "/reports/create?which=part",
	}
	return c.Render(200, "full_page_load", data)
}

func GetCreateReportPart(c echo.Context) error {
	locale := utils.GetLocale(c)
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "create_report", data)
}

func CreateReport(c echo.Context) error {
	locale := utils.GetLocale(c)
	var report model.Report
	description := c.FormValue("description")
	if description == "" {
		data := map[string]any{
			"locale": locale,
			"errors": map[string]string{
				"description": utils.Translate(locale, "report_create_empty_field"),
			},
		}
		return c.Render(400, "create_report", data)
	}
	report.Description = description
	if err := database.CreateReport(&report); err != nil {
		return err
	}
	data := map[string]any{
		"locale": locale,
	}
	return c.Render(200, "posts_main", data)
}

func ListReportsPaginated(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_MODERATOR.Level {
		return c.String(403, "Forbidden")
	}
	locale := utils.GetLocale(c)
	pageStr := c.QueryParam("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	reportsDB, err := database.GetReportsPaginated(page, 10)
	if err != nil {
		return c.String(500, "Internal server error")
	}
	reports := make([]map[string]any, len(reportsDB))
	for i, report := range reportsDB {
		desc := report.Description[:min(50, len(report.Description))]
		if len(report.Description) > 50 {
			desc += "..."
		}
		reports[i] = map[string]any{
			"id":          report.ID,
			"description": desc,
			"createdAt":   report.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	more := len(reports) == 10
	nextPageLoader := ""
	if more {
		nextPageLoader = fmt.Sprintf("/reports?page=%d", page+1)
	}
	data := map[string]any{
		"locale":   locale,
		"reports":  reports,
		"more":     more,
		"nextPage": nextPageLoader,
	}
	return c.Render(200, "reports_list", data)
}

func GetReport(c echo.Context) error {
	user, err := GetUserOfSession(c)
	if err != nil || user.Authority.Level < model.AUTH_MODERATOR.Level {
		return c.String(403, "Forbidden")
	}
	locale := utils.GetLocale(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.String(400, "Invalid report ID")
	}
	report, err := database.GetReportByID(id)
	if err != nil {
		return c.String(404, "Report not found")
	}
	data := map[string]any{
		"locale":      locale,
		"id":          report.ID,
		"description": report.Description,
		"createdAt":   report.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return c.Render(200, "report", data)
}
