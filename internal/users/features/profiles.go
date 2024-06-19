package features

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/pkg/errors"
	"github.com/uber/h3-go/v4"
	"go.uber.org/zap"
)

type UserProfilesService struct {
	userRepo data.IUserRepository
	config   *config.Config
	logger   *zap.Logger
}

func NewUserProfilesService(
	userRepo data.IUserRepository,
	cfg *config.Config,
	logger *zap.Logger,
) *UserProfilesService {
	return &UserProfilesService{
		userRepo: userRepo,
		config:   cfg,
		logger:   logger,
	}
}

// UpdateLocation updates the location of a user identified by their user ID, Should be an authenticated user.
// It retrieves the user from the repository using the provided user ID,
// calculates the H3 index based on the given longitude and latitude,
// and updates the user's location and H3 index in the repository.
// If the user is not found, it returns an ErrUserNotFound error.
// If the update fails, it returns an error wrapping the underlying error.
// It also logs a message indicating the successful update of the user's profile with the new location.
// Finally, it retrieves the updated user from the repository and returns it along with a nil error.
func (u *UserProfilesService) UpdateLocation(ctx context.Context, userId int64, long, lat float64) (models.User, error) {
	authUser, err := u.userRepo.FindOne(ctx, data.UserWithID(userId))
	if err != nil {
		return models.User{}, ErrUserNotFound
	}

	latLng := h3.NewLatLng(lat, long)
	const resolution = 9 // For good balance between precision (for a dating app usecase) and performance

	cell := h3.LatLngToCell(latLng, resolution)

	authUser.Latitude = lat
	authUser.Longitude = long
	authUser.H3Index = int64(cell)
	//user := models.User{
	//	ID:        authUser.ID,
	//	Longitude: long,
	//	Latitude:  lat,
	//	H3Index:   int64(cell),
	//}
	err = u.userRepo.Update(ctx, &authUser)
	if err != nil {
		return models.User{}, errors.Wrap(err, "update location failed")
	}
	u.logger.Info("Profile updated with location", zap.Int64("userID", authUser.ID))

	gotUser, _ := u.userRepo.FindOne(ctx, data.UserWithID(userId))
	return gotUser, nil
}
