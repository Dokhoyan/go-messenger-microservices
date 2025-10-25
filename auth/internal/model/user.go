package model

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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
	Username     string   `json:"username"`
	Role         UserRole `json:"role"`
	PasswordHash string   `json:"passwordhash"`
}

type User struct {
	ID        int64        `json:"id"`
	Info      UserAuthData `json:"info"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}
