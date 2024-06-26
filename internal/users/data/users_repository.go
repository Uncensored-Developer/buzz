package data

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/bun_mysql"
	"github.com/Uncensored-Developer/buzz/pkg/repository"
	"github.com/uber/h3-go/v4"
	"github.com/uptrace/bun"
	"time"
)

type IUserRepository interface {
	repository.IRepository[models.User]
	IncrementLikes(ctx context.Context, id int64, value int) error
}

func NewUserRepository(db bun.IDB) IUserRepository {
	return NewUserBunRepository(db)
}

type UserBunRepository struct {
	bun_mysql.BunRepository[models.User]
}

func NewUserBunRepository(db bun.IDB) *UserBunRepository {
	return &UserBunRepository{bun_mysql.BunRepository[models.User]{DB: db}}
}

func (b *UserBunRepository) IncrementLikes(ctx context.Context, id int64, value int) error {
	stmt := fmt.Sprintf("UPDATE users SET likes_count = likes_count + %d WHERE id = ?", value)
	_, err := b.DB.ExecContext(ctx, stmt, id)
	return err
}

func UserWithEmail(email string) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("email = ?", email)
	}
}

func UserWithID(id int64) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("id = ?", id)
	}
}

func UserWithEmailAndPassword(email, password string) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("email = ?", email).Where("password = ?", password)
	}
}

func UsersWithinDobRange(start, end time.Time) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("dob BETWEEN ? AND ?", start, end)
	}
}

func UsersWithGender(gender string) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("gender = ?", gender)
	}
}

func UsersExcludingID(id int64) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("id != ?", id)
	}
}

func UsersWithinH3Indexes(indexes []h3.Cell) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("h3_index IN (?)", bun.In(indexes))
	}
}
