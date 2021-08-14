package lib

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"userId"`
}

func ParseAuthToken(token string) (int, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(DefaultConfig.JwtSecret), nil
	})
	if err != nil || !jwtToken.Valid {
		return 0, fmt.Errorf("incorrect token")
	}

	if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
		return 0, fmt.Errorf("incorrect token")
	}

	claims := jwtToken.Claims.(*tokenClaims)
	if claims.UserID <= 0 {
		return 0, fmt.Errorf("incorrect token")
	}

	return claims.UserID, nil
}

func GenerateAuthToken(userID int) string {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
		UserID: userID,
	})

	tokenString, _ := jwtToken.SignedString([]byte(DefaultConfig.JwtSecret))
	return tokenString
}
