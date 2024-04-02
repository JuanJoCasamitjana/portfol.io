package routes

import (
	"github.com/JuanJoCasamitjana/portfol.io/internal/handlers"
	"github.com/labstack/echo/v4"
)

func setUpReportsRoutes(e *echo.Echo) {
	e.GET("/reports/create", handlers.GetCreateReport)
	e.POST("/reports/create", handlers.CreateReport)
	e.GET("/reports", handlers.ListReportsPaginated)
	e.GET("/reports/:id", handlers.GetReport)
}
