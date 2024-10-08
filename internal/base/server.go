package base

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database"
	"github.com/JuanJoCasamitjana/portfol.io/internal/routes"
	"github.com/JuanJoCasamitjana/portfol.io/internal/utils"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/natefinch/lumberjack.v2"
)

var port string
var SECRET string
var log_format = `{"time":${time_unix_milli},"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"}
`

func init() {
	//This is a basic configuration to launch the server on railway
	//with a dynamic port
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	env_port := os.Getenv("PORT")
	port = "8080"
	if env_port != "" {
		port = env_port
	}
	env_secret := os.Getenv("SECRET")
	SECRET = "SECRET"
	if env_secret != "" {
		SECRET = env_secret
	}

}

func SetUpAndRunServer() {
	e := echo.New()

	file_logger := &lumberjack.Logger{
		Filename:   "logs/portfol.io.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Compress:   true,
		LocalTime:  true,
	}
	mw := io.MultiWriter(os.Stdout, file_logger)
	log.SetOutput(mw)

	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: log_format,
			Output: mw,
		},
	),
	)
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(SECRET))))
	e.Renderer = NewTemplates()
	e.Static("/static", "web/static")
	routes.SetUpRoutes(e)
	shutdown := make(chan struct{})
	sysSignals := make(chan os.Signal, 1)
	signal.Notify(sysSignals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	go func() {
		e.Start("0.0.0.0:" + port)
	}()
	utils.Shutdown = &shutdown
	go func() {
		<-sysSignals
		utils.ShutDownSignal()
	}()
	<-shutdown
	defer cancel()
	err := e.Shutdown(ctx)
	if err != nil {
		fmt.Println("Error shutting down server", err)
	}
	fmt.Println("Server is shutting down")
	os.RemoveAll(database.Replicas)
	defer os.Exit(0)
}

// Templates is a custom renderer for echo
type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	funcMap := template.FuncMap{
		"Translate": utils.Translate,
	}
	return &Templates{
		templates: template.Must(template.New("").Funcs(funcMap).ParseGlob("./web/templates/*.html")),
	}
}
