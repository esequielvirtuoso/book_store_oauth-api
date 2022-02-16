// Package db holds the infrastructure layer to interact with databases repositories
package db

import (
	accesstoken "github.com/esequielvirtuoso/book_store_oauth-api/src/domain/access_token"
	"github.com/esequielvirtuoso/book_store_oauth-api/src/internal/infrastructure/clients/cassandra"
	restErrors "github.com/esequielvirtuoso/go_utils_lib/rest_errors"
	"github.com/gocql/gocql"
)

const (
	queryGetAccessToken    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?"
	queryInsertAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?)"
	queryUpdateExpires     = "UPDATE access_tokens SET expires=? WHERE access_token=?"
)

func NewService() DBRepository {
	return &dbRepository{}
}

type DBRepository interface {
	GetById(string) (*accesstoken.AccessToken, restErrors.RestErr)
	Create(accesstoken.AccessToken) restErrors.RestErr
	UpdateExpirationTime(accesstoken.AccessToken) restErrors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*accesstoken.AccessToken, restErrors.RestErr) {
	var result accesstoken.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(&result.AccessToken,
		&result.UserID,
		&result.AppClientID,
		&result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, restErrors.NewNotFoundError("no access token found with given id")
		}
		return nil, restErrors.NewInternalServerError(err.Error())
	}

	return &result, nil
}

func (r *dbRepository) Create(at accesstoken.AccessToken) restErrors.RestErr {
	if err := cassandra.GetSession().Query(queryInsertAccessToken,
		at.AccessToken,
		at.UserID,
		at.AppClientID,
		at.Expires).Exec(); err != nil {
		return restErrors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(at accesstoken.AccessToken) restErrors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateExpires, at.AccessToken).Exec(); err != nil {
		return restErrors.NewInternalServerError(err.Error())
	}
	return nil
}
