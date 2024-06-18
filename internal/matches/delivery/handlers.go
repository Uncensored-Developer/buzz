package delivery

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/matches/features"
	"github.com/Uncensored-Developer/buzz/internal/server/dto"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.uber.org/zap"
	"net/http"
)

func HandleUserSwipe(
	ctx context.Context,
	logger *zap.Logger,
	matchService *features.MatchService,
) http.Handler {
	type request struct {
		UserId int64  `json:"userId"`
		Action string `json:"action"`
	}
	type respResult struct {
		Matched   bool  `json:"matched"`
		MatchedId int64 `json:"matchedID,omitempty"`
	}
	type response struct {
		Results respResult `json:"results"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				return
			}
			logger.Info("Handling user swipe")

			authUser := r.Context().Value("user").(models.User)

			swipeInput, err := dto.DecodeValid[request](r)
			if err != nil {
				dto.SendErrorJsonResponse(w, logger, "Invalid request body format", http.StatusBadRequest)
				return
			}

			// Validate inputs and aggregate error messages if any
			err = validation.ValidateStruct(&swipeInput,
				validation.Field(&swipeInput.UserId, validation.Required),
				validation.Field(&swipeInput.Action, validation.Required),
				validation.Field(&swipeInput.Action, validation.In("YES", "NO")),
			)
			if err != nil {
				errMsgs := make(map[string]string)
				validationErrs := err.(validation.Errors)
				for name, err := range validationErrs {
					errMsgs[name] = err.Error()
				}
				dto.SendErrorJsonResponse(w, logger, errMsgs, http.StatusBadRequest)
				return
			}

			match, err := matchService.Swipe(ctx,
				authUser.ID, swipeInput.UserId, features.SwipeAction(swipeInput.Action))
			if err != nil {
				dto.SendErrorJsonResponse(w, logger, err.Error(), http.StatusBadRequest)
				return
			}

			successResp := response{
				Results: respResult{
					Matched: false,
				},
			}

			// Check if a match was returned
			if match.ID > 0 {
				// Update response with match ID
				successResp.Results.Matched = true
				successResp.Results.MatchedId = match.ID
			}
			err = dto.Encode[response](w, http.StatusOK, successResp)
			if err != nil {
				logger.Error("could not encode success response",
					zap.Error(err))
			}
			return
		},
	)
}
