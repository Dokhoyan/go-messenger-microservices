package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64          `db:"id"`
	Info      UserInfo       `db:""`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

type UserInfo struct {
	Name 	   string 		`db:"name"`
    Username   string 		`db:"username"`
    Email      string 		`db:"email"`
    Birth_date time.Time 	`db:"birth_date"`
    Avatar_url string 		`db:"avatar_url"`
}