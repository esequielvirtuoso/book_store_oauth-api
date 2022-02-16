// Package accesstoken holds all characteristics of access token domain or objects
package accesstoken

import (
	"fmt"
	"strings"
	"time"

	cryptoutils "github.com/esequielvirtuoso/go_utils_lib/crypto"
	restErrors "github.com/esequielvirtuoso/go_utils_lib/rest_errors"
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

func (at *AccessTokenRequest) Validate() *restErrors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		break
	case grantTypeClientCredentials:
		break
	default:
		return restErrors.NewBadRequestError("invalid grant_type")
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
func (at *AccessToken) Validate() *restErrors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return restErrors.NewBadRequestError("invalid access token id")
	}

	if at.UserID <= 0 {
		return restErrors.NewBadRequestError("invalid user id")
	}
	if at.AppClientID <= 0 {
		return restErrors.NewBadRequestError("invalid client id")
	}
	return nil
}

func (at *AccessToken) Generate() {
	at.AccessToken = cryptoutils.GetSha256(fmt.Sprintf("at-%d-%d-ran", at.UserID, at.Expires))
}
