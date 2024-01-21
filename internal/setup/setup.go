package setup

import (
	"os"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func SetupAndRun() {
	//Setup Database connection
	database.SetUpDB()

	//Setup logger config
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	app := echo.New()
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	//Add routing
	routes.SetUpRoutes(app)
	//I wonder if I can change the logger level during execution
	app.Logger.Info(app.Start(":8080"))
}
