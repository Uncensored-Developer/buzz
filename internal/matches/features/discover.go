package features

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type Gender string

const (
	MaleGender   Gender = "M"
	FemaleGender Gender = "F"
	OtherGender  Gender = "O"
)

type MatchFilter struct {
	MinAge int
	MaxAge int
	Gender Gender
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

// FetchPotentialMatches fetches potential matches for the given user ID and filters.
// It retrieves the authenticated user (userId), applies the filters, and returns a list of potential matches.
// If the user is not found, it returns ErrUserNotFound.
// The function uses the userRepo to retrieve matching users based on the provided filters.
func (d *DiscoverService) FetchPotentialMatches(ctx context.Context, userId int64, filters MatchFilter) ([]models.User, error) {
	authUser, err := d.userRepo.FindOne(ctx, data.UserWithID(userId))
	if err != nil {
		return nil, features.ErrUserNotFound
	}
	now := time.Now()

	var opts = []repository.SelectCriteria{
		data.UsersExcludingID(authUser.ID),
	}
	const maxAge int = 60
	const minAge int = 18

	if filters.MinAge > 0 && filters.MaxAge == 0 {
		// MinAge passed but no max age
		startDate := now.AddDate(-maxAge, 0, 0)
		endDate := now.AddDate(-filters.MinAge, 0, 0)
		opts = append(opts, data.UsersWithinDobRange(startDate, endDate))
	} else if filters.MinAge == 0 && filters.MaxAge > 0 {
		// MaxAge passed but no min age
		startDate := now.AddDate(-filters.MaxAge, 0, 0)
		endDate := now.AddDate(-minAge, 0, 0)
		opts = append(opts, data.UsersWithinDobRange(startDate, endDate))
	} else if filters.MinAge > 0 && filters.MaxAge > 0 {
		// MinAge and MaxAge passed
		startDate := now.AddDate(-filters.MaxAge, 0, 0)
		endDate := now.AddDate(-filters.MinAge, 0, 0)
		opts = append(opts, data.UsersWithinDobRange(startDate, endDate))
	}

	if filters.Gender != "" {
		opts = append(opts, data.UsersWithGender(string(filters.Gender)))
	}

	users, err := d.userRepo.FindAll(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "fetch users failed")
	}
	return users, nil
}
