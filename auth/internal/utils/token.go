package utils

import (
	"errors"
	"time"

	
	"github.com/Dokhoyan/go-messenger-microservices/auth/internal/model"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userAuthData model.UserAuthData, secretKey []byte, duration time.Duration) (string, error){
	claims:=model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: userAuthData.Username,
		Role: userAuthData.Role,
	}

	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}


func VerifyToken(tokenHash string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenHash,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("wrong type signing method")
			}

			return secretKey, nil
		})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}