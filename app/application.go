package app

import (
	"github.com/dula0/bookstore_oauth_api/clients/cassandra"
	"github.com/dula0/bookstore_oauth_api/http"
	"github.com/dula0/bookstore_oauth_api/repository/db"
	"github.com/dula0/bookstore_oauth_api/src/domain/access_token"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine = gin.Default()
)

func StartApp() {
	session, dbErr := cassandra.GetSession()
	if dbErr != nil {
		panic(dbErr)
	}
	defer session.Close()

	// access token handler
	atHandler := http.NewHandler(access_token.NewService(db.NewRepo()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)

	router.Run(":8080")
}
