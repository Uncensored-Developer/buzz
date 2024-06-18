package data

import (
	"github.com/Uncensored-Developer/buzz/internal/matches/models"
	"github.com/Uncensored-Developer/buzz/pkg/bun_mysql"
	"github.com/Uncensored-Developer/buzz/pkg/repository"
	"github.com/uptrace/bun"
)

type IMatchesRepository interface {
	repository.IRepository[models.Match]
}

func NewMatchesRepository(db bun.IDB) IMatchesRepository {
	return bun_mysql.NewBunRepository[models.Match](db)
}

func MatchWithID(id int64) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("id = ?", id)
	}
}

func MatchWithUserOneID(id int64) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("user_one_id = ?", id)
	}
}

func MatchWithUserTwoID(id int64) repository.SelectCriteria {
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Where("user_two_id = ?", id)
	}
}
