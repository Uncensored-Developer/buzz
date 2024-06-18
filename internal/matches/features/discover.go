package features

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type MatchFilter struct {
	MinAge int
	MaxAge int
	Gender string
}

type DiscoverService struct {
	userRepo data.IUserRepository
	config   *config.Config
	logger   *zap.Logger
}

func NewDiscoverService(
	userRepo data.IUserRepository,
	cfg *config.Config,
	logger *zap.Logger,
) *DiscoverService {
	return &DiscoverService{
		userRepo: userRepo,
		config:   cfg,
		logger:   logger,
	}
}

func (d *DiscoverService) FetchPotentialMatches(ctx context.Context, userId int64, filters MatchFilter) ([]models.User, error) {
	authUser, err := d.userRepo.FindOne(ctx, data.UserWithID(userId))
	if err != nil {
		return nil, features.ErrUserNotFound
	}

	if filters.Gender == "" {
		if authUser.Gender == "M" {
			filters.Gender = "F"
		} else {
			filters.Gender = "M"
		}
	}
	fmt.Printf("%+v\n", filters)
	now := time.Now()

	var opts = []repository.SelectCriteria{
		data.UsersExcludingID(authUser.ID),
	}
	const maxAge int = 60
	const minAge int = 18

	if filters.MinAge > 0 && filters.MaxAge == 0 {
		// MinAge passed but no max age
		minDate := now.AddDate(-filters.MinAge, 0, 0)
		maxDate := now.AddDate(-maxAge-1, 0, 0).AddDate(0, 0, 1)
		opts = append(opts, data.UsersWithinDobRange(minDate, maxDate))
	} else if filters.MinAge == 0 && filters.MaxAge > 0 {

		// MaxAge passed but no min age
		minDate := now.AddDate(-minAge, 0, 0)
		maxDate := now.AddDate(-filters.MaxAge-1, 0, 0).AddDate(0, 0, 1)
		opts = append(opts, data.UsersWithinDobRange(minDate, maxDate))
	} else {
		opts = append(opts, data.UsersWithGender(filters.Gender))
	}

	users, err := d.userRepo.FindAll(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "fetch users failed")
	}
	return users, nil
}
