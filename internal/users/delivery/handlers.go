package delivery

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/server/dto"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/brianvoe/gofakeit/v7"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// HandleUserCreate handles the creation of a new user via an empty HTTP POST request.
// It creates/sign's up a new user by generating random user profile information
func HandleUserCreate(
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
	authService *features.AuthenticationService,
) http.Handler {

	type userResponse struct {
		Id       int64  `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Age      int    `json:"age"`
	}

	// HTTP response type for user signup
	type response struct {
		Result userResponse `json:"result"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				return
			}
			logger.Info("Handling user signup")

			// generate random user details for signup
			email := gofakeit.Email()
			name := gofakeit.Name()
			password := cfg.FakeUserPassword
			gender := strings.ToUpper(string([]rune(gofakeit.Gender())[0]))
			dob := gofakeit.PastDate()

			user, err := authService.SignUp(ctx, dob, name, email, password, gender)
			if err != nil {
				dto.SendErrorJsonResponse(w, logger, err.Error())
			} else {
				age := time.Now().Year() - user.Dob.Year()
				res := response{
					Result: userResponse{
						Id:       user.ID,
						Email:    user.Email,
						Password: user.Password,
						Name:     user.Name,
						Gender:   user.Gender,
						Age:      age,
					},
				}
				err := dto.Encode[response](w, 201, res)
				if err != nil {
					logger.Error("could not encode success response",
						zap.Error(err))
				}
			}
		},
	)
}

func HandleUserLogin(
	ctx context.Context,
	logger *zap.Logger,
	cfg *config.Config,
	authService *features.AuthenticationService,
) http.Handler {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type successResponse struct {
		Token string `json:"token"`
	}

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodPost {
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
				return
			}
			logger.Info("Handling user login")

			loginInput, err := dto.DecodeValid[request](r)
			if err != nil {
				dto.SendErrorJsonResponse(w, logger, "Invalid request body format")
				return
			}

			// Validate inputs and aggregate error messages if any
			err = validation.ValidateStruct(&loginInput,
				validation.Field(&loginInput.Email, validation.Required),
				validation.Field(&loginInput.Password, validation.Required),
			)
			if err != nil {
				errMsgs := make(map[string]string)
				validationErrs := err.(validation.Errors)
				for name, err := range validationErrs {
					errMsgs[name] = err.Error()
				}
				dto.SendErrorJsonResponse(w, logger, errMsgs)
				return
			}

			token, err := authService.Login(ctx, loginInput.Email, loginInput.Password)
			if err != nil {
				dto.SendErrorJsonResponse(w, logger, err.Error())
				return
			}

			response := successResponse{token}
			err = dto.Encode[successResponse](w, http.StatusOK, response)
			if err != nil {
				logger.Error("could not encode success response",
					zap.Error(err))
			}
			return
		},
	)
}
