package db

import (
	"github.com/dula0/bookstore_oauth_api/clients/cassandra"
	"github.com/dula0/bookstore_oauth_api/src/domain/access_token"
	"github.com/dula0/bookstore_users_api/utils/errors"
)

// cassandra query
const (
	getAccessTokenQuery    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	createAccessTokenQuery = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	UpdateExpirationQuery = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

func NewRepo() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	session, err := cassandra.GetSession()
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	defer session.Close()

	var result access_token.AccessToken
	if err := session.Query(getAccessTokenQuery, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires,
	); err != nil {
		if err.Error() == "not found" {
			return nil, errors.NotFoundError("no access token found with given id")
		}
		return nil, errors.InternalServerError(err.Error())
	}

	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) *errors.RestErr {
	session, err := cassandra.GetSession()
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	defer session.Close()

	if err := session.Query(createAccessTokenQuery,
		at.AccessToken,
		at.UserId,
		at.ClientId,
		at.Expires,
	).Exec(); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	session, err := cassandra.GetSession()
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	defer session.Close()

	if err := session.Query(createAccessTokenQuery,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}
