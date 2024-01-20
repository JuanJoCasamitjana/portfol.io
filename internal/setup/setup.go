package setup

import (
	"os"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var config = fiber.Config{}

func SetupAndRun() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	app := fiber.New(config)
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	app.Listen(":3000")
}
