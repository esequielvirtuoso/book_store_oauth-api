// Package accesstoken holds all characteristics of access token domain or objects
package accesstoken

import (
	"time"
)

const (
	expirationTime = 24
)

// AccessToken defines access token characteristics
type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserID      int64  `json:"user_id`
	AppClientID int64  `json:"app_client_id"`
	Expires     int64  `json:"expires"`
}

// GetNewAccessToken return a new access token
func GetNewAccessToken() AccessToken {
	return AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

// IsExpired validate if the access token is expired or not
func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}
