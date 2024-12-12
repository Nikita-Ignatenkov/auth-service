package handlers

import (
	"auth-service/models"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUser(c *gin.Context) {
	// Логирование начала обработки запроса на получение пользователя
	logrus.Info("Handling get user request")

	userID := c.Param("id")
	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		// Логирование ошибки поиска пользователя
		logrus.Errorf("User not found: %v", err)
		return
	}

	// Логирование успешного получения пользователя
	c.JSON(http.StatusOK, user)
	logrus.Info("User retrieved successfully")
}

func ChangePassword(c *gin.Context) {
	// Логирование начала обработки запроса на изменение пароля
	logrus.Info("Handling change password request")

	userID := c.MustGet("user_id").(uint)
	var input struct {
		CurrentPass string `json:"current_password" binding:"required"`
		NewPass     string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// Логирование ошибки валидации данных
		logrus.Errorf("Validation error: %v", err)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		// Логирование ошибки поиска пользователя
		logrus.Errorf("User not found: %v", err)
		return
	}

	if !models.CheckPasswordHash(input.CurrentPass, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect current password"})
		// Логирование ошибки аутентификации
		logrus.Errorf("Incorrect current password")
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
