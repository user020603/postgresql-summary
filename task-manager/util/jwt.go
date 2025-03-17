package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWTToken(userID uint, username string) (string, int64, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret_key" 
	}

	var expHours int64 = 24
	expStr := os.Getenv("JWT_EXPIRATION")
	if expStr != "" {
		if val, err := strconv.ParseInt(expStr, 10, 64); err == nil {
			expHours = val
		}
	}

	expiration := time.Now().Add(time.Hour * time.Duration(expHours)).Unix()

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "task-manager",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiration - time.Now().Unix(), nil
}

func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret_key" 
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}