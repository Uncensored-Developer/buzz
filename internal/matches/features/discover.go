package features

import (
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"go.uber.org/zap"
)

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
