package utils

import (
	"fmt"
	"os"


	"github.com/golang-jwt/jwt/v5"
)

func GenToken(payload map[string]interface{}) (string, error) {
	secret := os.Getenv("SECRET")
	if payload == nil {
		return "", fmt.Errorf("payload doesn't exist")
	}

	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyToken(tokenStr string) (jwt.MapClaims, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token is empty")
	}
	secret := os.Getenv("SECRET")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
