// Package db holds the infrastructure layer to interact with databases repositories
package db

import (
	accesstoken "github.com/esequielvirtuoso/book_store_oauth-api/src/domain/access_token"
	"github.com/esequielvirtuoso/book_store_oauth-api/src/domain/access_token/utils/errors"
)

func NewService() DBRepository {
	return &dbRepository{}
}

type DBRepository interface {
	GetById(string) (*accesstoken.AccessToken, *errors.RestErr)
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*accesstoken.AccessToken, *errors.RestErr) {
	return nil, errors.NewInternalServerError("database connection not implemented yet")
}
