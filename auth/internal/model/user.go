package model

import "github.com/dgrijalva/jwt-go"

type UserRole int8

const (
	UNKNOWN UserRole = iota
	USER
	ADMIN
)

const (
	IDFieldCode       = "id"
	UserNameFieldCode = "username"
)

// UserClaims - параметры для JWT токена
type UserClaims struct {
	jwt.StandardClaims
	Username string
	Role     UserRole
}

type UserAuthData struct{
	Username string
	Role     UserRole
	Password string
}