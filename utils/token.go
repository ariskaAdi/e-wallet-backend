package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string, verified bool, secret string) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"id":  id,
		"verified": verified,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (id string,role string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "","", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id = fmt.Sprintf("%v", claims["id"])
		role = fmt.Sprintf("%v", claims["role"])
		return 
	}

	err = fmt.Errorf("unable to extract claims")
	return 
}