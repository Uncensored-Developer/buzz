package models

import (
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/uptrace/bun"
	"time"
)

type Match struct {
	bun.BaseModel `bun:"table:matches"`

	ID        int64     `bun:"id,pk,autoincrement"`
	UserOneID int64     `bun:"user_one_id,notnull"`
	UserTwoID int64     `bun:"user_two_id,notnull"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`

	UserOne *models.User `bun:"rel:belongs-to"`
	UserTwo *models.User `bun:"rel:belongs-to"`
}
