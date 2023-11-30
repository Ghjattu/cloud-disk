package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key = "secret"

func GenerateToken(userID int64, name string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func ValidateToken(tokenString string) (int64, string, error) {
	// If the token is empty, return an error.
	if tokenString == "" {
		return -1, "", fmt.Errorf("empty token")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(key), nil
	})

	if clams, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(clams["user_id"].(float64))
		name := clams["name"].(string)

		return userID, name, nil
	}

	return -1, "", err
}
