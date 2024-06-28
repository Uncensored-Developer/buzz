package features

import (
	"context"
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/datastore"
	data2 "github.com/Uncensored-Developer/buzz/internal/matches/data"
	models2 "github.com/Uncensored-Developer/buzz/internal/matches/models"
	"github.com/Uncensored-Developer/buzz/internal/users/data"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/Uncensored-Developer/buzz/pkg/repository"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type SwipeAction string

const (
	YesAction SwipeAction = "YES"
	NoAction  SwipeAction = "NO"
)

type MatchService struct {
	userRepo     data.IUserRepository
	cacheManager repository.ISimpleCacheManager
	uow          datastore.IUnitOfWork
	config       *config.Config
	logger       *zap.Logger
}

func NewMatchService(
	userRepo data.IUserRepository,
	cacheManager repository.ISimpleCacheManager,
	uow datastore.IUnitOfWork,
	cfg *config.Config,
	logger *zap.Logger,
) *MatchService {
	return &MatchService{
		userRepo:     userRepo,
		cacheManager: cacheManager,
		uow:          uow,
		config:       cfg,
		logger:       logger,
	}
}

// Swipe performs the swipe action between two users.
// It checks if the swiper user and swiped user exist in the user repository.
// If the action is "YesAction", it checks if there is a previous like from the swiped user.
// If there is no previous like, it saves the swipe action to the cache database.
// If there is a previous like, it saves the match to the matches repository and deletes the swipe action from the cache database.
// For other swipe actions like "NoAction" it does not perform any additional action.
// Parameters:
// - ctx: The context.Context for the operation.
// - swiperUserID: The ID of the user performing the swipe action, SHOULD BE THE AUTHENTICATED USER
// - swipedUserID: The ID of the user being swiped.
// - action: The swipe action to perform.
// Returns:
// - models.Match: The match object if a match is found, otherwise an empty match object.
// - error: An error if any occurred during the operation.
func (m *MatchService) Swipe(
	ctx context.Context,
	swiperUserID, swipedUserID int64,
	action SwipeAction,
) (models2.Match, error) {
	authUser, err := m.userRepo.FindOne(ctx, data.UserWithID(swiperUserID))
	if err != nil {
		return models2.Match{}, features.ErrUserNotFound
	}
	swipedUser, err := m.userRepo.FindOne(ctx, data.UserWithID(swipedUserID))
	if err != nil {
		return models2.Match{}, features.ErrUserNotFound
	}
	// This is the format the key for all save swipe account should follow
	// e.g. user 1 likes user 2 = `1.YES.2`
	// e.g. user 3 passes user 5 = `3.NO.5`
	swipeActionTemplate := "%d." + string(action) + ".%d"

	var gotMatch models2.Match
	err = m.uow.Do(ctx, func(store datastore.IUnitOfWorkDatastore) error {
		if action == YesAction {
			// Check if swipedUser has previously liked user's profile
			getKey := fmt.Sprintf(swipeActionTemplate, swipedUserID, swiperUserID)
			_, err := m.cacheManager.Get(ctx, getKey)
			if err != nil {
				// No match yet, just save to cache database
				saveKey := fmt.Sprintf(swipeActionTemplate, swiperUserID, swipedUserID)
				cacheDuration := time.Hour * 24 * 30 // 30 days
				err = m.cacheManager.Set(ctx, saveKey, string(action), cacheDuration)
				if err != nil {
					m.logger.Error("cache save failed", zap.Error(err))
					return errors.Wrap(err, "cache save failed")
				}
				m.logger.Info("swipe action saved",
					zap.Int64("userId", swiperUserID),
					zap.String("action", string(action)),
					zap.Int64("swipedUserId", swipedUserID),
				)
			} else {
				// Match found, save to matches repository
				match := &models2.Match{
					UserOneID: authUser.ID,
					UserTwoID: swipedUser.ID,
				}
				err := store.MatchesRepository().Save(ctx, match)
				if err != nil {
					m.logger.Error("match save failed", zap.Error(err))
					return errors.Wrap(err, "match save failed")
				}

				gM, _ := store.MatchesRepository().FindOne(ctx,
					data2.MatchWithUserOneID(match.UserOneID),
					data2.MatchWithUserTwoID(match.UserTwoID),
				)
				m.logger.Info("match occured",
					zap.Int64("matchID", gotMatch.ID),
				)

				// Delete swipe action from cache database
				err = m.cacheManager.Delete(ctx, getKey)
				if err != nil {
					m.logger.Error("swipe action delete failed", zap.Error(err))
					return errors.Wrap(err, "swipe action delete failed")
				}
				gotMatch = gM
			}

			// Increment likes count for swiped user
			err = store.UsersRepository().IncrementLikes(ctx, swipedUserID, 1)
			if err != nil {
				m.logger.Error("increment likes failed", zap.Error(err))
				return errors.Wrap(err, "increment likes failed")
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return models2.Match{}, errors.Wrap(err, "handle swipe failed")
	}
	// Handle other swipe accounts like NO, SUPER LIKE etc.
	return gotMatch, nil
}
