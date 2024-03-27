package utils

import (
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var (
	FromEmail         string
	FromEmailPassword []byte
	SmtpHost          string
	SmtpPort          string
)

func init() {
	godotenv.Load()
	FromEmail = os.Getenv("FROM_EMAIL")
	FromEmailPassword = []byte(os.Getenv("FROM_EMAIL_PASSWORD"))
	SmtpHost = os.Getenv("SMTP_HOST")
}

func GetLocale(c echo.Context) string {
	lang := c.Request().Header.Get("Accept-Language")
	if lang == "" {
		return "en"
	}
	lang = strings.Split(lang, "-")[0]
	return lang
}

func SendEmailNotification(to []string, body []byte) error {
	if len(to) == 0 {
		return nil
	}
	auth := smtp.PlainAuth("", FromEmail, string(FromEmailPassword), SmtpHost)
	return smtp.SendMail(SmtpHost+":"+SmtpPort, auth, FromEmail, to, body)
}
