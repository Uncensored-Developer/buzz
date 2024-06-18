package features

import (
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/test/e2e"
	"github.com/Uncensored-Developer/buzz/internal/users/models"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
	"time"
)

type userLoginTestSuite struct {
	suite.Suite
	e2e.TestServerSuite
	currentUser *models.User
}

func (u *userLoginTestSuite) SetupSuite() {
	err := u.StartUp()
	if err != nil {
		u.Logger.Fatal("Test setup failed", zap.Error(err))
	}
	user, err := e2e.CreateUser(u.Ctx, time.Time{}, "", u.Config.FakeUserPassword, "M")
	if err != nil {
		u.Logger.Fatal("create test user failed", zap.Error(err))
	}
	u.currentUser = user
}

func (u *userLoginTestSuite) TearDownSuite() {
	err := u.Shutdown()
	if err != nil {
		u.Logger.Fatal("Test shutdown failed", zap.Error(err))
	}
}

func (u *userLoginTestSuite) TestLoginWithEmptyFields() {
	url := fmt.Sprintf("%s/login", u.ServerURL)

	type fieldResp struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type errorResponse struct {
		Error fieldResp `json:"error"`
	}

	const expectedStatus = 400
	expectedLoginResponse := errorResponse{
		Error: fieldResp{
			Email:    "cannot be blank",
			Password: "cannot be blank",
		},
	}
	loginInput := fieldResp{
		Email:    "",
		Password: "",
	}
	client := resty.New()
	res, err := client.R().SetBody(loginInput).SetError(&errorResponse{}).Post(url)
	if err != nil {
		u.Logger.Error("req client error", zap.Error(err))
	}
	u.Require().NoError(err)

	u.Assert().Equal(expectedStatus, res.StatusCode())
	u.Assert().Equal(&expectedLoginResponse, res.Error())
}

func (u *userLoginTestSuite) TestLoginWithInvalidCredentials() {
	url := fmt.Sprintf("%s/login", u.ServerURL)

	type errorResponse struct {
		Error string `json:"error"`
	}
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	const expectedStatus = 400
	expectedLoginResponse := errorResponse{
		Error: "Email or password incorrect",
	}

	loginInput := loginReq{
		Email:    u.currentUser.Email,
		Password: "wrongPassword",
	}

	client := resty.New()
	res, err := client.R().SetBody(loginInput).SetError(&errorResponse{}).Post(url)
	if err != nil {
		u.Logger.Error("req client error", zap.Error(err))
	}
	u.Require().NoError(err)

	u.Assert().Equal(expectedStatus, res.StatusCode())
	u.Assert().Equal(&expectedLoginResponse, res.Error())
}

func (u *userLoginTestSuite) TestLoginWithValidCredentials() {
	url := fmt.Sprintf("%s/login", u.ServerURL)

	type successResponse struct {
		Token string `json:"token"`
	}
	type loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	const expectedStatus = 200

	loginInput := loginReq{
		Email:    u.currentUser.Email,
		Password: u.Config.FakeUserPassword,
	}

	client := resty.New()
	res, err := client.R().SetBody(loginInput).SetResult(&successResponse{}).Post(url)
	if err != nil {
		u.Logger.Error("req client error", zap.Error(err))
	}
	u.Require().NoError(err)

	u.Assert().Equal(expectedStatus, res.StatusCode())
	u.Assert().NotEmpty(res.Result())
}

func TestUserLoginE2e(t *testing.T) {
	suite.Run(t, new(userLoginTestSuite))
}
