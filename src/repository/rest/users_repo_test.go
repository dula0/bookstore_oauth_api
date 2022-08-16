package rest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("starting test cases")
	rest.StartMockupServer() //comment out the flag.Parse() in mock.go or you'll get an error
	os.Exit(m.Run())
}

func TestLoginUserTimeout(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          fmt.Sprintf("%v/users/login", BASEURL),
		ReqBody:      `{"email":"test@gmail.com","password":"password"`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repo := usersRepo{}
	user, err := repo.LoginUser("test@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restclient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          fmt.Sprintf("%v/users/login", BASEURL),
		ReqBody:      `{"email":"test@gmail.com","password":"password"`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid login credentials", "status": "404", "error": "not_found"}`,
	})

	repo := usersRepo{}
	user, err := repo.LoginUser("test@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          fmt.Sprintf("%v/users/login", BASEURL),
		ReqBody:      `{"email":"test@gmail.com","password":"password"`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid login credentials", "status": 404, "error": "not_found"}`,
	})

	repo := usersRepo{}
	user, err := repo.LoginUser("test@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          fmt.Sprintf("%v/users/login", BASEURL),
		ReqBody:      `{"email":"test@gmail.com","password":"password"`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{ "id": "1", "first_name": "test", "last_name": "name", "email": "test@gmail.com"}`,
	})

	repo := usersRepo{}
	user, err := repo.LoginUser("test@gmail.com", "password")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to decode users login response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          fmt.Sprintf("%v/users/login", BASEURL),
		ReqBody:      `{"email":"test@gmail.com","password":"password"`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{ "id": 1, "first_name": "test", "last_name": "name", "email": "test@gmail.com"}`,
	})

	repo := usersRepo{}
	user, err := repo.LoginUser("test@gmail.com", "password")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "test", user.FirstName)
	assert.EqualValues(t, "name", user.LastName)
	assert.EqualValues(t, "test@gmail.com", user.Email)
}
