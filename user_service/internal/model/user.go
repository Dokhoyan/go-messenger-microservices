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
	Name 	   string 		`db:"name" json:"name"`
    Username   string 		`db:"username" json:"username"`
    Email      string 		`db:"email" json:"email"`
    Birth_date time.Time 	`db:"birth_date" json:"birth_date"`
    Avatar_url string 		`db:"avatar_url" json:"avatar_url"`
}

