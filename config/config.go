package config

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Config struct {
		Port      string
		DBDriver  string
		DBSource  string
		SMTPHost  string
		SMTPPort  int
		SMTPUser  string
		SMTPPass  string
		SecretKey string
	}
)

func LoadConfig() {
	// Логирование начала загрузки конфигурации
	logrus.Info("Loading configuration from .env file")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Config.Port = os.Getenv("PORT")
	Config.DBDriver = os.Getenv("DB_DRIVER")
	Config.DBSource = os.Getenv("DB_SOURCE")
	Config.SMTPHost = os.Getenv("SMTP_HOST")
	Config.SMTPPort, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))
	Config.SMTPUser = os.Getenv("SMTP_USER")
	Config.SMTPPass = os.Getenv("SMTP_PASSWORD")
	Config.SecretKey = os.Getenv("SECRET_KEY")

	// Логирование завершения загрузки конфигурации
	logrus.Info("Configuration loaded successfully")
}
