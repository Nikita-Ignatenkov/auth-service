package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	MaxLoginAttempts = 5
	BanDuration      = 5 * time.Minute
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func BruteForceProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Логирование начала проверки попыток входа
		logrus.Info("Checking login attempts")

		email := c.PostForm("email")

		// Получаем количество попыток входа из Redis
		ctx := context.Background()
		attempts, err := redisClient.Get(ctx, email).Result()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check login attempts"})
			c.Abort()
			// Логирование ошибки проверки попыток входа
			logrus.Errorf("Failed to check login attempts: %v", err)
			return
		}

		if attempts == "" {
			attempts = "0"
		}

		// Преобразуем количество попыток в число
		attemptCount, err := strconv.Atoi(attempts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse login attempts"})
			c.Abort()
			// Логирование ошибки парсинга попыток входа
			logrus.Errorf("Failed to parse login attempts: %v", err)
			return
		}

		if attemptCount >= MaxLoginAttempts {
			// Если превышено максимальное количество попыток, блокируем пользователя
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many login attempts. Please try again later."})
			c.Abort()
			// Логирование блокировки пользователя
			logrus.Warn("Too many login attempts")
			return
		}

		// Увеличиваем количество попыток
		attemptCount++
		err = redisClient.Set(ctx, email, strconv.Itoa(attemptCount), BanDuration).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increment login attempts"})
			c.Abort()
			// Логирование ошибки увеличения попыток входа
			logrus.Errorf("Failed to increment login attempts: %v", err)
			return
		}

		// Логирование успешной проверки попыток входа
		logrus.Info("Login attempts checked successfully")
		c.Next()
	}
}
