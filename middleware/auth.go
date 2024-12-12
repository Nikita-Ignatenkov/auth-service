package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Логирование начала проверки токена
		logrus.Info("Authenticating request")

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			// Логирование ошибки отсутствия токена
			logrus.Warn("Token missing")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(c.GetString("secret_key")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			// Логирование ошибки невалидного токена
			logrus.Warn("Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			// Логирование ошибки невалидных утверждений токена
			logrus.Warn("Invalid token claims")
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id claim"})
			c.Abort()
			// Логирование ошибки невалидного утверждения user_id
			logrus.Warn("Invalid user_id claim")
			return
		}

		c.Set("user_id", uint(userID))
		c.Next()

		// Логирование успешной аутентификации
		logrus.Info("Request authenticated successfully")
	}
}
