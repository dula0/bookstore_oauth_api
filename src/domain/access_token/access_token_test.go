package access_token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenConstants(t *testing.T) {
	assert.EqualValues(t, 24, expirationTime, "expiration time should be 24 hours")
}

func TestGetNewAccessToken(t *testing.T) {
	at := GetNewAccessToken()

	assert.False(t, at.IsExpired(), "new access token should not be nil")

	assert.EqualValues(t, "", at.AccessToken, "new access token shouldn't have token id")

	assert.True(t, at.UserId == 0, "new access token should not have a user id")
}

func TestAccessToken_IsExpired(t *testing.T) {
	at := AccessToken{}

	assert.True(t, at.IsExpired(), "empty access token should be expired by default")

	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "an access token set three hours ago should not be expired")
}
