package main

import (
	"auth-service/config"
	"auth-service/database"
	"auth-service/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"os"
)

func main() {
	// Логирование начала выполнения приложения
	logrus.Info("Starting the auth-service application")

	config.LoadConfig()
	db := database.ConnectDB()

	r := gin.Default()

	// Логирование начала подключения статических файлов
	logrus.Info("Serving static files")

	// Подключение статических файлов
	r.Use(static.Serve("/static", static.LocalFile("./static", true)))

	// Логирование завершения подключения статических файлов
	logrus.Info("Static files served successfully")

	// Подключение сессий
	store := cookie.NewStore([]byte(config.Config.SecretKey))
	r.Use(sessions.Sessions("mysession", store))

	// Подключение мидлварей
	r.Use(database.Middleware(db))

	// Роуты
	routes.SetupRoutes(r)

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Настройка logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	// Запись логов в файл
	f, err := os.OpenFile("auth-service.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	logrus.SetOutput(f)

	// Логирование перед запуском сервера
	logrus.Info("Starting server on :" + config.Config.Port)

	r.Run(":" + config.Config.Port)

	// Логирование после завершения выполнения приложения
	logrus.Info("Stopping the auth-service application")
}
