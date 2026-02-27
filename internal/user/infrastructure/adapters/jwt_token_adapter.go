package adapters

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTTokenAdapter struct {
	secret string
}

func NewJWTTokenAdapter() *JWTTokenAdapter {
	return &JWTTokenAdapter{secret: os.Getenv("JWT_SECRET")}
}

func (j *JWTTokenAdapter) GenerateToken(userId int, email string, name string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"name":    name,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTTokenAdapter) ValidateToken(tokenString string) (bool, map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return false, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil, fmt.Errorf("claims inválidos")
	}

	result := make(map[string]interface{})
	for k, v := range claims {
		result[k] = v
	}

	return true, result, nil
}
