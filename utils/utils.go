package utils

import (
	"auth-service/config"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/smtp"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func SendRegistrationConfirmation(email string) error {
	msg := "Subject: Registration Confirmation\n\nThank you for registering with our service."

	auth := smtp.PlainAuth("", config.Config.SMTPUser, config.Config.SMTPPass, config.Config.SMTPHost)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", config.Config.SMTPHost, config.Config.SMTPPort), auth, "your_email@example.com", []string{email}, []byte(msg))
	if err != nil {
		log.Printf("Failed to send registration confirmation email: %v", err)
		return err
	}

	return nil
}

func SendResetPasswordLink(email string) error {
	token, err := GenerateJWT(0)
	if err != nil {
		log.Printf("Failed to generate reset password token: %v", err)
		return err
	}

	msg := fmt.Sprintf("Subject: Reset Password\n\nClick the following link to reset your password: http://localhost:8080/reset-password?token=%s", token)

	auth := smtp.PlainAuth("", config.Config.SMTPUser, config.Config.SMTPPass, config.Config.SMTPHost)
	err = smtp.SendMail(fmt.Sprintf("%s:%d", config.Config.SMTPHost, config.Config.SMTPPort), auth, "your_email@example.com", []string{email}, []byte(msg))
	if err != nil {
		log.Printf("Failed to send reset password link: %v", err)
		return err
	}

	return nil
}

func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.SecretKey))
}

func ParseJWT(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.Config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, false
	}

	return token.Claims.(jwt.MapClaims), true
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
