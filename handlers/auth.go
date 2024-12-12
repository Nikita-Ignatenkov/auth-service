package handlers

import (
	"auth-service/models"
	"auth-service/utils"
	"github.com/gin-contrib/sessions"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	// Логирование начала обработки запроса на регистрацию
	logrus.Info("Handling registration request")

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// Логирование ошибки валидации данных
		logrus.Errorf("Validation error: %v", err)
		return
	}

	input.Password, _ = utils.HashPassword(input.Password)

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		// Логирование ошибки создания пользователя
		logrus.Errorf("Failed to create user: %v", err)
		return
	}

	// Отправка подтверждения регистрации через email
	err := utils.SendRegistrationConfirmation(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send registration confirmation email"})
		// Логирование ошибки отправки письма
		logrus.Errorf("Failed to send registration confirmation email: %v", err)
		return
	}

	// Генерация токена после успешной регистрации
	token, err := utils.GenerateJWT(input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token after registration"})
		// Логирование ошибки генерации токена
		logrus.Errorf("Failed to generate token after registration: %v", err)
		return
	}

	// Логирование успешной регистрации
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "token": token})
	logrus.Info("User registered successfully")

	// Переадресация на страницу входа или профиля
	c.Redirect(http.StatusSeeOther, "/login")
}

func Login(c *gin.Context) {
	// Логирование начала обработки запроса на вход
	logrus.Info("Handling login request")

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// Логирование ошибки валидации данных
		logrus.Errorf("Validation error: %v", err)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		// Логирование ошибки аутентификации
		logrus.Errorf("Authentication failed: %v", err)
		return
	}

	if !models.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		// Логирование ошибки аутентификации
		logrus.Errorf("Authentication failed: invalid password")
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		// Логирование ошибки генерации токена
		logrus.Errorf("Failed to generate token: %v", err)
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	// Логирование успешной аутентификации
	c.JSON(http.StatusOK, gin.H{"token": token})
	logrus.Info("User authenticated successfully")
}

func ResetPasswordRequest(c *gin.Context) {
	// Логирование начала обработки запроса на восстановление пароля
	logrus.Info("Handling reset password request")

	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// Логирование ошибки валидации данных
		logrus.Errorf("Validation error: %v", err)
		return
	}

	user := models.User{}
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		// Логирование ошибки поиска пользователя
		logrus.Errorf("User not found: %v", err)
		return
	}

	// Отправка ссылки для сброса пароля через email
	err := utils.SendResetPasswordLink(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset password link"})
		// Логирование ошибки отправки письма
		logrus.Errorf("Failed to send reset password link: %v", err)
		return
	}

	// Логирование успешной отправки ссылки на восстановление пароля
	c.JSON(http.StatusOK, gin.H{"message": "Reset password link sent"})
	logrus.Info("Reset password link sent successfully")
}

func ResetPassword(c *gin.Context) {
	// Логирование начала обработки запроса на сброс пароля
	logrus.Info("Handling reset password request")

	var input struct {
		Token   string `json:"token" binding:"required"`
		NewPass string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// Логирование ошибки валидации данных
		logrus.Errorf("Validation error: %v", err)
		return
	}

	claims, ok := utils.ParseJWT(input.Token)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// Логирование ошибки проверки токена
		logrus.Errorf("Invalid token")
		return
	}

	userID := claims["user_id"].(float64)
	user := models.User{}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		// Логирование ошибки поиска пользователя
		logrus.Errorf("User not found: %v", err)
		return
	}

	newHashedPassword, _ := models.HashPassword(input.NewPass)
	user.Password = newHashedPassword

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		// Логирование ошибки обновления пароля
		logrus.Errorf("Failed to update password: %v", err)
		return
	}

	// Логирование успешного обновления пароля
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	logrus.Info("Password updated successfully")
}
