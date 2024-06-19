package delivery

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/server/dto"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func HandleUpdateProfileLocation(
	ctx context.Context,
	logger *zap.Logger,
	profileService *features.UserProfilesService,
) http.Handler {

	// HTTP request type for profile update
	type updateRequest struct {
		Longitude float64
		Latitude  float64
	}

	type userResponse struct {
		Id        int64   `json:"id"`
		Email     string  `json:"email"`
		Password  string  `json:"password"`
		Name      string  `json:"name"`
		Gender    string  `json:"gender"`
		Age       int     `json:"age"`
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
		H3Index   int64   `json:"h3_index"`
	}

	// HTTP response type for profile update
	type response struct {
		Result userResponse `json:"result"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				return
			}
			logger.Info("Handling update profile location")
			authUser := r.Context().Value("user").(models.User)

			updateInput, err := dto.DecodeValid[updateRequest](r)
			if err != nil {
				dto.SendErrorJsonResponse(w, logger, "Invalid request body format", http.StatusBadRequest)
				return
			}

			// Validate inputs and aggregate error messages if any
			err = validation.ValidateStruct(&updateInput,
				validation.Field(&updateInput.Longitude, validation.Required),
				validation.Field(&updateInput.Latitude, validation.Required),
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

			user, err := profileService.UpdateLocation(ctx, authUser.ID, updateInput.Longitude, updateInput.Latitude)
			if err != nil {
				dto.SendErrorJsonResponse(w, logger, err.Error(), http.StatusBadRequest)
				return
			}
			age := time.Now().Year() - user.Dob.Year()
			res := response{
				Result: userResponse{
					Id:        user.ID,
					Email:     user.Email,
					Password:  user.Password,
					Name:      user.Name,
					Gender:    user.Gender,
					Age:       age,
					Longitude: user.Longitude,
					Latitude:  user.Latitude,
					H3Index:   user.H3Index,
				},
			}
			err = dto.Encode[response](w, http.StatusOK, res)
			if err != nil {
				logger.Error("could not encode success response",
					zap.Error(err))
			}
			return
		},
	)
}
