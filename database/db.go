package database

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	// Логирование начала подключения к базе данных
	logrus.Info("Connecting to the database")

	dsn := os.Getenv("DB_SOURCE")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("failed to connect database: %v", err)
	}

	// Логирование завершения подключения к базе данных
	logrus.Info("Database connection established")

	return db
}

func Middleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
