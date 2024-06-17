package delivery

import (
	"context"
	"github.com/Uncensored-Developer/buzz/internal/server/dto"
	"github.com/Uncensored-Developer/buzz/internal/users/features"
	"github.com/Uncensored-Developer/buzz/pkg/config"
	"github.com/brianvoe/gofakeit/v7"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

func HandleUserSignUp(
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
				errRes := dto.ErrorResponse{
					Error: err.Error(),
				}
				err := dto.Encode[dto.ErrorResponse](w, r, http.StatusInternalServerError, errRes)
				if err != nil {
					logger.Error("could not encode error response",
						zap.Error(err))
				}
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
				err := dto.Encode[response](w, r, 201, res)
				if err != nil {
					logger.Error("could not encode success response",
						zap.Error(err))
				}
			}
		},
	)
}
