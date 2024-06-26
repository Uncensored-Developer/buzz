package features

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/authentication"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/hash"
	"github.com/pkg/errors"
	"github.com/uber/h3-go/v4"
	"go.uber.org/zap"
	"time"
)

var ErrEmailTaken = errors.New("Email already taken")
var ErrInvalidLoginCred = errors.New("Email or password incorrect")

type AuthenticationService struct {
	hasher       hash.IStringHasher
	tokenManager authentication.ITokenManager
	userRepo     data.IUserRepository
	config       *config.Config
	logger       *zap.Logger
}

func NewAuthenticationService(
	hasher hash.IStringHasher,
	tokenManager authentication.ITokenManager,
	userRepo data.IUserRepository,
	cfg *config.Config,
	logger *zap.Logger,
) *AuthenticationService {
	return &AuthenticationService{
		hasher:       hasher,
		tokenManager: tokenManager,
		userRepo:     userRepo,
		config:       cfg,
		logger:       logger,
	}
}

// SignUp registers a new user in the system.
// It takes the user's date of birth (dob), gender, name, email, and password.
// It returns the created User object and an error, if any.
func (a *AuthenticationService) SignUp(
	ctx context.Context,
	dob time.Time,
	lat, long float64,
	name, email, password, gender string) (models.User, error) {

	_, err := a.userRepo.FindOne(ctx, data.UserWithEmail(email))
	if err == nil {
		return models.User{}, ErrEmailTaken
	}

	hashedPassword, err := a.hasher.Hash(password)
	if err != nil {
		return models.User{}, errors.Wrap(err, "password hash failed")
	}

	latLng := h3.NewLatLng(lat, long)
	cell := h3.LatLngToCell(latLng, a.config.H3Resolution)

	newUser := &models.User{
		Email:     email,
		Password:  hashedPassword,
		Name:      name,
		Gender:    gender,
		Dob:       dob,
		Latitude:  lat,
		Longitude: long,
		H3Index:   int64(cell),
	}
	err = a.userRepo.Save(ctx, newUser)
	if err != nil {
		a.logger.Error("User save operation failed",
			zap.Any("user", newUser),
		)
		return models.User{}, errors.Wrap(err, "user save failed")
	}
	user, _ := a.userRepo.FindOne(ctx, data.UserWithEmail(newUser.Email))

	a.logger.Info("User successfully created",
		zap.String("email", user.Email),
		zap.Int64("id", user.ID),
	)
	return user, nil
}

func (a *AuthenticationService) Login(ctx context.Context, email, password string) (string, error) {
	hashPassword, err := a.hasher.Hash(password)
	if err != nil {
		return "", errors.Wrap(err, "password hash failed")
	}
	user, err := a.userRepo.FindOne(ctx, data.UserWithEmailAndPassword(email, hashPassword))
	if err != nil {
		return "", ErrInvalidLoginCred
	}
	sixMonthsDuration := time.Hour * 24 * 30 * 6
	token, err := a.tokenManager.NewToken(user.ID, sixMonthsDuration)
	if err != nil {
		return "", errors.Wrap(err, "token encoding failed")
	}
	return token, nil
}

func (a *AuthenticationService) GetUserFromToken(ctx context.Context, token string) (models.User, error) {
	userId, err := a.tokenManager.Parse(token)
	if err != nil {
		return models.User{}, errors.Wrap(err, "token parse failed")
	}

	user, err := a.userRepo.FindOne(ctx, data.UserWithID(userId))
	if err != nil {
		return models.User{}, errors.Wrap(err, "user not found")
	}
	return user, nil
}
