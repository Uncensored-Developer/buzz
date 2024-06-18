package features

import (
	"fmt"
	"github.com/Uncensored-Developer/buzz/internal/test/e2e"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

type createUserE2eTestSuite struct {
	suite.Suite
	e2e.TestServerSuite
}

func (c *createUserE2eTestSuite) SetupSuite() {
	err := c.StartUp()
	if err != nil {
		c.Logger.Fatal("Test setup failed", zap.Error(err))
	}
}

func (c *createUserE2eTestSuite) TearDownSuite() {
	err := c.Shutdown()
	if err != nil {
		c.Logger.Fatal("Test shutdown failed", zap.Error(err))
	}
}

func (c *createUserE2eTestSuite) TestCreateUserRouteOnlyAllowPostRequest() {
	url := fmt.Sprintf("%s/user/create", c.ServerURL)

	testCases := map[string]int{
		http.MethodPut:    405,
		http.MethodDelete: 405,
		http.MethodGet:    405,
		http.MethodPost:   201,
	}
	for method, expectedStatus := range testCases {
		c.T().Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, url, nil)
			c.Require().NoError(err)

			res, err := http.DefaultClient.Do(req)
			c.Require().NoError(err)

			c.Assert().Equal(res.StatusCode, expectedStatus)
		})
	}
}

func (c *createUserE2eTestSuite) TestCreateRandomUser() {
	url := fmt.Sprintf("%s/user/create", c.ServerURL)

	type userResponse struct {
		Id       int64  `json:"id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Age      int    `json:"age"`
	}

	// HTTP response type for user signup
	type successResponse struct {
		Result userResponse `json:"result"`
	}

	var userResp successResponse

	const expectedStatus = 201

	client := resty.New()
	res, err := client.R().SetResult(&userResp).Post(url)
	if err != nil {
		c.Logger.Error("req client error", zap.Error(err))
	}
	c.Require().NoError(err)

	c.Assert().Equal(res.StatusCode(), expectedStatus)
	c.Assert().NotEmpty(userResp)
}

func TestCreateUserE2e(t *testing.T) {
	suite.Run(t, new(createUserE2eTestSuite))
}
