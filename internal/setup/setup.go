package setup

import (
	"html/template"
	"io"
	"os"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/routes"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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
	app.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	app.Renderer = NewTemplates()
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			contentType := c.Request().Header.Get("Content-Type")
			logger.Debug().
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("method", c.Request().Method).
				Str("type", contentType).
				Err(v.Error).
				Msg("request")

			return nil
		},
	}))
	app.Static("/static", "web/static")
	//Add routing
	routes.SetUpRoutes(app)
	//I wonder if I can change the logger level during execution
	app.Start(":8080")
}

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("./web/templates/*.html")),
	}
}
