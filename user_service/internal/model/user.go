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

type User struct {
	ID        int64          `db:"id"`
	Info      UserInfo       `db:""`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
	Password  string         `db:"password"`
}

type UserInfo struct {
	Name 	   string 		`db:"name"`
    Username   string 		`db:"username"`
    Email      string 		`db:"email"`
    Birth_date time.Time 	`db:"birth_date"`
    Avatar_url string 		`db:"avatar_url"`
	Role       UserRole     `db:"role"`
}

type UserCreate struct {
	Info     UserInfo `db:""`
	Password string   `db:"password"`
}

// UserUpdate - DTO для обновления пользователя
type UserUpdate struct {
	ID   int64    `db:"id"`
	Info UserInfo `db:""`
}

// UserClaims - параметры для JWT токена
type UserClaims struct {
	jwt.StandardClaims
	Username string
	Role     UserRole
}

func (r UserRole) String() string {
	switch r {
	case USER:
		return "USER"
	case ADMIN:
		return "ADMIN"
	default:
		return "UNKNOWN"
	}
}