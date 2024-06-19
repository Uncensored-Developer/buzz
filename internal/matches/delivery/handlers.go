package delivery

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/matches/features"
	"github.com/Uncensored-Developer/buzz/internal/server/dto"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/Uncensored-Developer/buzz/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strconv"
	"time"
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

func HandleFetchPotentialMatches(
	ctx context.Context,
	logger *zap.Logger,
	discService *features.DiscoverService,
) http.Handler {

	type userResp struct {
		Id             int64  `json:"id"`
		Name           string `json:"name"`
		Gender         string `json:"gender"`
		Age            int    `json:"age"`
		DistanceFromMe int    `json:"distanceFromMe"`
	}

	type successResp struct {
		Results []userResp `json:"results"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				return
			}
			logger.Info("Handling fetch potential matches")
			authUser := r.Context().Value("user").(models.User)

			searchGender := ""
			if authUser.Gender == "M" {
				searchGender = "F"
			} else {
				searchGender = "M"
			}

			minAge := 18
			maxAge := 60
			ageRangeStr := r.URL.Query().Get("age_range")
			if ageRangeStr != "" {
				ageRangePattern := `^(\d+)-(\d+)$`
				re := regexp.MustCompile(ageRangePattern)

				// Check if ageRangeStr matches pattern
				matches := re.FindStringSubmatch(ageRangeStr)
				if matches == nil {
					var msg = map[string]string{"age_range": "Invalid age range format"}
					dto.SendErrorJsonResponse[map[string]string](w, logger, msg, http.StatusBadRequest)
					return
				}

				start, err1 := strconv.Atoi(matches[1])
				minAge = start
				end, err2 := strconv.Atoi(matches[2])
				maxAge = end
				if err1 != nil || err2 != nil {
					var msg = map[string]string{"age_range": "Invalid range value"}
					dto.SendErrorJsonResponse[map[string]string](w, logger, msg, http.StatusBadRequest)
					return
				}
			}

			gender := r.URL.Query().Get("gender")
			if gender != "" {
				if !isValidGender(gender) {
					var msg = map[string]string{"gender": "Invalid gender value"}
					dto.SendErrorJsonResponse[map[string]string](w, logger, msg, http.StatusBadRequest)
					return
				}
				searchGender = gender
			}

			distance := r.URL.Query().Get("distance_from")
			radius, err := strconv.Atoi(distance)
			if err != nil {
				var msg = map[string]string{"gender": "Invalid distance_from value"}
				dto.SendErrorJsonResponse[map[string]string](w, logger, msg, http.StatusBadRequest)
				return
			}

			filters := features.MatchFilter{
				MinAge: minAge,
				MaxAge: maxAge,
				Gender: features.Gender(searchGender),
				Radius: radius,
			}
			users, err := discService.FetchPotentialMatches(ctx, authUser.ID, filters)
			if err != nil {
				dto.SendErrorJsonResponse[string](w, logger, err.Error(), http.StatusBadRequest)
				return
			}

			var usersResp []userResp
			for _, user := range users {
				userResponse := userResp{
					Id:     user.ID,
					Name:   user.Name,
					Gender: user.Gender,
					Age:    time.Now().Year() - user.Dob.Year(),
					DistanceFromMe: int(utils.DistanceBetween(
						authUser.Latitude, authUser.Longitude,
						user.Latitude, user.Longitude)),
				}
				usersResp = append(usersResp, userResponse)
			}
			resp := successResp{usersResp}

			err = dto.Encode[successResp](w, http.StatusOK, resp)
			if err != nil {
				logger.Error("could not encode success response",
					zap.Error(err))
			}
			return
		},
	)
}

func isValidGender(gender string) bool {
	return gender == "M" || gender == "F" || gender == "O"
}
