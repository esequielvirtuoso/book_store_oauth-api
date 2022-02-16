// Package accesstoken test holds all tests over access token domain
package accesstoken

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAcessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "expiration time should be 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	accessToken := GetNewAccessToken(0)

	assert.False(t, accessToken.IsExpired(), "brand new access token should not be expired")

	assert.Empty(t, accessToken.AccessToken, "new access token should not have defined access token id")

	assert.Zero(t, accessToken.UserID, "new access token should not have an associated user id")
}

func TestIsExpired(t *testing.T) {
	accessToken := AccessToken{}
	assert.True(t, accessToken.IsExpired(), "empty access token should be expired by default")

	accessToken.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, accessToken.IsExpired(), "access token expiring three hours from now should not be expired")
}
