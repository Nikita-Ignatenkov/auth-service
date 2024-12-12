package handlers

import (
	"auth-service/models"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ManageUsers(c *gin.Context) {
	logrus.Info("Handling manage users request")

	db := c.MustGet("db").(*gorm.DB)
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		logrus.Errorf("Failed to fetch users: %v", err) // Убедитесь, что err является значением ошибки
		return
	}

	c.JSON(http.StatusOK, users)
	logrus.Info("Users retrieved successfully")
}

func AssignRole(c *gin.Context) {
	logrus.Info("Handling assign role request")

	var input struct {
		UserID uint `json:"user_id" binding:"required"`
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Errorf("Validation error: %v", err) // Убедитесь, что err является значением ошибки
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	var role models.Role

	if err := db.Preload("Roles").First(&user, input.UserID).Error; err != nil { // Убедитесь, что здесь используется Error
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		logrus.Errorf("User not found: %v", err) // Убедитесь, что err является значением ошибки
		return
	}

	if err := db.First(&role, input.RoleID).Error; err != nil { // Убедитесь, что здесь используется Error
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		logrus.Errorf("Role not found: %v", err) // Убедитесь, что err является значением ошибки
		return
	}

	if err := db.Model(&user).Association("Roles").Append(&role).Error; err != nil { // Убедитесь, что здесь используется Error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
	logrus.Info("Role assigned successfully")
}
