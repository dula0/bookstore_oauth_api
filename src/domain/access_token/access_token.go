package access_token

// domain definition

import (
	"fmt"
	"strings"
	"time"

	"github.com/dula0/bookstore_users_api/utils/crypto_utils"
	"github.com/dula0/bookstore_users_api/utils/errors"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grandTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	Username string `json:"username"`
	Password string `json:"password"`

	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *errors.RestErr {

	if at.GrantType != grantTypePassword || at.GrantType != grandTypeClientCredentials {
		return errors.BadRequestError("invalid grant type")
	}
	
	/*
	switch at.GrantType {
	case grantTypePassword:
		break

	case grandTypeClientCredentials:
		break

	default:
		return errors.BadRequestError("invalid grant_type parameter")
	}

	
	*/
	return nil
}

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

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}
