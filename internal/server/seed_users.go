package server

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/bun_mysql"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/db"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/pkg/errors"
	"github.com/uber/h3-go/v4"
	"go.uber.org/zap"
	"time"
)

func PreLoadUsers(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.Logger,
) {
	testUsers := []models.User{
		{
			Email:     "test.user@buzz.com",
			Password:  "77546b6a714655717144705547696e414a71485acbfdac6008f9cab4083784cbd1874f76618d2a97",
			Name:      gofakeit.Name(),
			Gender:    "F",
			Dob:       time.Date(1999, 3, 24, 0, 0, 0, 0, time.UTC),
			Longitude: 0.5026768,
			Latitude:  51.2725887,
		},
		{
			Email:     "user22@buzz.com",
			Password:  "77546b6a714655717144705547696e414a71485acbfdac6008f9cab4083784cbd1874f76618d2a97",
			Name:      gofakeit.Name(),
			Gender:    "M",
			Dob:       time.Date(2001, 3, 24, 0, 0, 0, 0, time.UTC),
			Latitude:  50.96284649,
			Longitude: -0.12981616,
		},
		{
			Email:     "user32@buzz.com",
			Password:  "77546b6a714655717144705547696e414a71485acbfdac6008f9cab4083784cbd1874f76618d2a97",
			Name:      gofakeit.Name(),
			Gender:    "M",
			Dob:       time.Date(1995, 3, 24, 0, 0, 0, 0, time.UTC),
			Latitude:  51.03052403,
			Longitude: 0.18169958,
		},
		{
			Email:     "user42@buzz.com",
			Password:  "77546b6a714655717144705547696e414a71485acbfdac6008f9cab4083784cbd1874f76618d2a97",
			Name:      gofakeit.Name(),
			Gender:    "F",
			Dob:       time.Date(2004, 2, 24, 0, 0, 0, 0, time.UTC),
			Longitude: 0.62093072,
			Latitude:  51.31488722,
		},
		{
			Email:     "user52@buzz.com",
			Password:  "77546b6a714655717144705547696e414a71485acbfdac6008f9cab4083784cbd1874f76618d2a97",
			Name:      gofakeit.Name(),
			Gender:    "O",
			Dob:       time.Date(1988, 3, 24, 0, 0, 0, 0, time.UTC),
			Latitude:  6.524379,
			Longitude: 3.379206,
		},
	}
	bunDb, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("could make db connection", zap.Error(err))
		return
	}
	userRepo := data.NewUserRepository(bunDb)
	for _, user := range testUsers {
		_, err := userRepo.FindOne(ctx, data.UserWithEmail(user.Email))
		if errors.Is(err, bun_mysql.ErrRowNotFound) {
			err = userRepo.Save(ctx, &user)
			if err != nil {
				logger.Error("could not create users", zap.Error(err))
			}
		}
	}

	for _, user := range testUsers {
		gotUser, _ := userRepo.FindOne(ctx, data.UserWithEmail(user.Email))

		latLng := h3.NewLatLng(user.Latitude, user.Latitude)
		cell := h3.LatLngToCell(latLng, cfg.H3Resolution)
		gotUser.H3Index = int64(cell)
		err := userRepo.Update(ctx, &gotUser)
		if err != nil {
			logger.Error("could not update H3 index", zap.Error(err))
		}
	}
	logger.Info("Preloaded users")
}
