package app

import (
	"github.com/dula0/bookstore_oauth_api/src/clients/cassandra"
	"github.com/dula0/bookstore_oauth_api/src/http"
	"github.com/dula0/bookstore_oauth_api/src/repository/db"
	"github.com/dula0/bookstore_oauth_api/src/repository/rest"
	"github.com/dula0/bookstore_oauth_api/src/services/access_token"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine = gin.Default()
)

func StartApp() {
	session := cassandra.GetSession()
	defer session.Close()

	// access token handler
	atHandler := http.NewHandler(access_token.NewService(rest.NewRepo(), db.NewRepo()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8080")
}
