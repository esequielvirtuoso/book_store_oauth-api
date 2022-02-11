// Package accesstoken holds all characteristics of access token domain or objects
package accesstoken

import (
	"fmt"
	"strings"
	"time"

	cryptoUtils "github.com/esequielvirtuoso/book_store_oauth-api/src/utils/crypto_utils"
	"github.com/esequielvirtuoso/book_store_oauth-api/src/utils/errors"
)

const (
	expirationTime             = 30
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

// AccessTokenRequest defines access token body request
type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant_type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client_credentials grant type
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *errors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredentials:
		break
	default:
		return errors.NewBadRequestError("invalid grant_type")
	}
	// TODO: Validate parameters for each grant_type
	return nil
}

// AccessToken defines access token characteristics
type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id"`
	AppClientID int64  `json:"app_client_id,omitempty"`
	Expires     int64  `json:"expires"`
}

// GetNewAccessToken return a new access token
func GetNewAccessToken(userID int64) AccessToken {
	return AccessToken{
		UserID:  userID,
		Expires: time.Now().UTC().Add(expirationTime * time.Minute).Unix(),
	}
}

// IsExpired validate if the access token is expired or not
func (at *AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

// Validate function verify if access token is valid or not
func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestError("invalid access token id")
	}

	if at.UserID <= 0 {
		return errors.NewBadRequestError("invalid user id")
	}
	if at.AppClientID <= 0 {
		return errors.NewBadRequestError("invalid client id")
	}
	return nil
}

func (at *AccessToken) Generate() {
	at.AccessToken = cryptoUtils.GetSha256(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}
