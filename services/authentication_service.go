package services

import (
	"arithmetic-calculator/config"
	"arithmetic-calculator/models"
	"arithmetic-calculator/utils"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Authenticate(username, password string) (string, error) {
	var user models.User
	if err := config.DB.Where("username = ? ", username).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}
	if utils.CheckPasswordHash(user.Password, password) {
		return "", errors.New("incorrect password")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
