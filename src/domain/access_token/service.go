// Package accesstoken service holds all business logic to handle access token operations
package accesstoken

import (
	"strings"

	"github.com/esequielvirtuoso/book_store_oauth-api/src/domain/access_token/utils/errors"
)

type Repository interface {
	GetById(string) (*AccessToken, *errors.RestErr)
}

type Service interface {
	GetById(string) (*AccessToken, *errors.RestErr)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (s *service) GetById(accessTokenID string) (*AccessToken, *errors.RestErr) {
	accessTokenID = strings.TrimSpace(accessTokenID)
	if len(accessTokenID) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := s.repository.GetById(accessTokenID)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}
