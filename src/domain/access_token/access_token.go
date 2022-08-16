package access_token

// domain definition

import (
	"strings"
	"time"

	"github.com/dula0/bookstore_users_api/utils/errors"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) Validate() *errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.BadRequestError("invalid access token id")
	}

	if at.UserId <= 0 {
		return errors.BadRequestError("invalid user id")

	}
	if at.ClientId <= 0 {
		return errors.BadRequestError("invalid client id")

	}
	if at.Expires <= 0 {
		return errors.BadRequestError("invalid expiration date")

	}

	return nil
}

func GetNewAccessToken() *AccessToken {
	return &AccessToken{
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}
