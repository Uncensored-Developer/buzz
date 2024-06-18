package features

import (
	"context"
	data2 "github.com/Uncensored-Developer/buzz/internal/matches/data"
	models2 "github.com/Uncensored-Developer/buzz/internal/matches/models"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"go.uber.org/zap"
)

type SwipeAction string

const (
	Like SwipeAction = "LIKE"
	Pass SwipeAction = "PASS"
)

type MatchService struct {
	userRepo    data.IUserRepository
	matchesRepo data2.IMatchesRepository
	config      *config.Config
	logger      *zap.Logger
}

func NewMatchService(
	userRepo data.IUserRepository,
	cfg *config.Config,
	logger *zap.Logger,
) *MatchService {
	return &MatchService{
		userRepo: userRepo,
		config:   cfg,
		logger:   logger,
	}
}

func (m *MatchService) Swipe(ctx context.Context, user models.User, action SwipeAction) (models2.Match, error) {
	return models2.Match{}, nil
}
