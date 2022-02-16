// Package accesstoken service holds all business logic to handle access token operations
package accesstoken

import (
	"strings"

	accesstoken "github.com/esequielvirtuoso/book_store_oauth-api/src/domain/access_token"
	"github.com/esequielvirtuoso/book_store_oauth-api/src/internal/infrastructure/repository/db"
	"github.com/esequielvirtuoso/book_store_oauth-api/src/internal/infrastructure/repository/rest"
	restErrors "github.com/esequielvirtuoso/go_utils_lib/rest_errors"
)

type Repository interface {
	GetById(string) (*accesstoken.AccessToken, *restErrors.RestErr)
	Create(accesstoken.AccessToken) *restErrors.RestErr
	UpdateExpirationTime(accesstoken.AccessToken) *restErrors.RestErr
}

type Service interface {
	GetById(string) (*accesstoken.AccessToken, *restErrors.RestErr)
	Create(accesstoken.AccessTokenRequest) (*accesstoken.AccessToken, *restErrors.RestErr)
	UpdateExpirationTime(accesstoken.AccessToken) *restErrors.RestErr
}

type service struct {
	repository db.DBRepository
	restClient rest.UsersClient
}

func NewService(userRepo rest.UsersClient, dbRepo db.DBRepository) Service {
	return &service{
		repository: dbRepo,
		restClient: userRepo,
	}
}

func (s *service) GetById(accessTokenID string) (*accesstoken.AccessToken, *restErrors.RestErr) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return nil, restErrors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := s.repository.GetById(accessTokenID)
	if err != nil {
		return nil, err
	}

	if accessToken.IsExpired() {
		return nil, restErrors.NewUnauthorized("access token id has expired")
	}
	return accessToken, nil
}

func (s *service) Create(request accesstoken.AccessTokenRequest) (*accesstoken.AccessToken, *restErrors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Authenticate the user against the Users API:
	user, err := s.restClient.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token:
	at := accesstoken.GetNewAccessToken(user.ID)
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.repository.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at accesstoken.AccessToken) *restErrors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}

	return s.repository.UpdateExpirationTime(at)
}
