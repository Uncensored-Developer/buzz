package features

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/authentication"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/hash"
	"github.com/pkg/errors"
	"time"
)

var ErrEmailTaken = errors.New("Email already taken")

type AuthenticationService struct {
	hasher       hash.IStringHasher
	tokenManager authentication.ITokenManager
	userRepo     data.IUserRepository
	config       *config.Config
}

func NewAuthenticationService(
	hasher hash.IStringHasher,
	tokenManager authentication.ITokenManager,
	userRepo data.IUserRepository,
	cfg *config.Config,
) *AuthenticationService {
	return &AuthenticationService{
		hasher:       hasher,
		tokenManager: tokenManager,
		userRepo:     userRepo,
		config:       cfg,
	}
}

func (a *AuthenticationService) SignUp(
	ctx context.Context,
	dob time.Time,
	gender, name, email, password string) (models.User, error) {

	_, err := a.userRepo.FindOne(ctx, data.UserWithEmail(email))
	if err == nil {
		return models.User{}, ErrEmailTaken
	}

	hashedPassword, err := a.hasher.Hash(password)
	if err != nil {
		return models.User{}, errors.Wrap(err, "password hash failed")
	}

	newUser := &models.User{
		Email:    email,
		Password: hashedPassword,
		Name:     name,
		Gender:   gender,
		Dob:      dob,
	}
	err = a.userRepo.Save(ctx, newUser)
	if err != nil {
		return models.User{}, errors.Wrap(err, "user save failed")
	}

	user, _ := a.userRepo.FindOne(ctx, data.UserWithEmail(newUser.Email))
	return user, nil
}
