// Package rest holds the logic to make http calls to other Rest APIs
package rest

import (
	"encoding/json"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"

	"github.com/esequielvirtuoso/book_store_oauth-api/src/domain/users"
	env "github.com/esequielvirtuoso/go_utils_lib/envs"
	restErrors "github.com/esequielvirtuoso/go_utils_lib/rest_errors"
)

const (
	envUserAPIURL = "USER_API_URL"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: getUserAPIURL(),
		Timeout: 100 * time.Millisecond,
	}
)

type UsersClient interface {
	LoginUser(string, string) (*users.User, restErrors.RestErr)
}

type usersClient struct{}

func NewClient() UsersClient {
	env.CheckRequired(envUserAPIURL)
	return &usersClient{}
}

func (c *usersClient) LoginUser(email string, password string) (*users.User, restErrors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, restErrors.NewInternalServerError("invalid rest client response when trying to login user", nil)
	}
	if response.StatusCode > 299 {
		apiErr, err := restErrors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, restErrors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, apiErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, restErrors.NewInternalServerError("error when trying to unmarshal users response", err)
	}
	return &user, nil
}

func getUserAPIURL() string {
	return env.GetString(envUserAPIURL)
}
