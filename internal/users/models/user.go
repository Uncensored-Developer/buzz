package models

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Email     string    `bun:"email,notnull"`
	Password  string    `bun:"password,notnull"`
	Name      string    `bun:"name,notnull"`
	Gender    string    `bun:"gender,notnull"`
	Dob       time.Time `bun:"dob,notnull,type:date"`
	Longitude float64   `bun:"longitude"`
	Latitude  float64   `bun:"latitude"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}
