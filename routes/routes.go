package routes

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Логирование начала настройки роутов
	logrus.Info("Setting up routes")

	authGroup := router.Group("/api")
	{
		authGroup.Use(middleware.Logging(), middleware.AuthMiddleware())
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", middleware.BruteForceProtection(), handlers.Login)
		authGroup.POST("/reset-password-request", handlers.ResetPasswordRequest)
		authGroup.POST("/reset-password", handlers.ResetPassword)
		authGroup.GET("/user/:id", handlers.GetUser)
		authGroup.PUT("/change-password", handlers.ChangePassword)
		authGroup.GET("/manage/users", handlers.ManageUsers)
		authGroup.POST("/assign-role", handlers.AssignRole)
	}

	// Роуты для обслуживания HTML страниц
	router.GET("/", func(c *gin.Context) {
		// Логирование начала обработки запроса на главную страницу
		logrus.Info("Handling root route request")
		c.File("./static/index.html")
	})

	router.GET("/register", func(c *gin.Context) {
		// Логирование начала обработки запроса на регистрацию
		logrus.Info("Handling register route request")
		c.File("./templates/register.html")
	})

	router.GET("/login", func(c *gin.Context) {
		// Логирование начала обработки запроса на вход
		logrus.Info("Handling login route request")
		c.File("./templates/login.html")
	})

	router.GET("/profile", func(c *gin.Context) {
		// Логирование начала обработки запроса на профиль
		logrus.Info("Handling profile route request")
		c.File("./templates/profile.html")
	})

	router.GET("/admin", func(c *gin.Context) {
		// Логирование начала обработки запроса на административную панель
		logrus.Info("Handling admin route request")
		c.File("./templates/admin.html")
	})

	// Логирование завершения настройки роутов
	logrus.Info("Routes setup completed")
}
