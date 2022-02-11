// Package rest holds the logic to make http calls to other Rest APIs
package rest

import (
	"encoding/json"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"

	"github.com/esequielvirtuoso/book_store_oauth-api/src/domain/users"
	"github.com/esequielvirtuoso/book_store_oauth-api/src/utils/errors"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:5001",
		Timeout: 100 * time.Millisecond,
	}
)

type UsersClient interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersClient struct{}

func NewClient() UsersClient {
	return &usersClient{}
}

func (c *usersClient) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServerError("invalid rest client response when trying to login user")
	}
	if response.StatusCode > 299 {
		var restErr errors.RestErr
		if err := json.Unmarshal(response.Bytes(), &restErr); err != nil {
			return nil, errors.NewInternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.NewInternalServerError("error when trying to unmarshal users response")
	}
	return &user, nil
}
