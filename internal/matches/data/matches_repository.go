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
