package rest

import (
	"encoding/json"
	"os"
	"time"

	"github.com/dula0/bookstore_oauth_api/src/domain/users"
	"github.com/dula0/bookstore_users_api/utils/errors"
	"github.com/mercadolibre/golang-restclient/rest"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: BASEURL,
		Timeout: 150 * time.Millisecond,
	}

	BASEURL = os.Getenv("BASEURL")
)

type RestUsersRepo interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepo struct{}

func NewRepo() RestUsersRepo {
	return &usersRepo{}
}

func (r *usersRepo) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response := usersRestClient.Post("/users/login", request)

	if response == nil || response.Response == nil {
		return nil, errors.InternalServerError("invalid restclient response when trying to login user")
	}

	if response.StatusCode > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, errors.InternalServerError("invalid error interface when trying to login user")
		}
		return nil, &restErr
	}

	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, errors.InternalServerError("error when trying to decode users login response")
	}
	return &user, nil
}
