package db

import (
	"github.com/dula0/bookstore_oauth_api/src/clients/cassandra"
	"github.com/dula0/bookstore_oauth_api/src/domain/access_token"
	"github.com/dula0/bookstore_users_api/utils/errors"
)

// cassandra query
const (
	getAccessTokenQuery    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	createAccessTokenQuery = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	UpdateExpirationQuery  = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
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

	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(getAccessTokenQuery, id).Scan(
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

	if err := cassandra.GetSession().Query(createAccessTokenQuery,
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

	if err := cassandra.GetSession().Query(UpdateExpirationQuery,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}
