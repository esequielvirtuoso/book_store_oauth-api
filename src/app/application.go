// Package app starts the application and map the HTTP routes.
package app

import (
	accesstoken "github.com/esequielvirtuoso/book_store_oauth-api/src/domain/access_token"
	"github.com/esequielvirtuoso/book_store_oauth-api/src/internal/infrastructure/http"
	"github.com/esequielvirtuoso/book_store_oauth-api/src/internal/infrastructure/repository/db"
	"github.com/gin-gonic/gin"
)

var (
	// NOTE: This is the only layer we are defining and using the HTTP server
	router = gin.Default()
)

// StartApplication attempts to map the API routes
func StartApplication() {
	// create a repository service
	repositoryService := db.NewService()

	// create a access token service and inject the repository service into this
	accessTokenService := accesstoken.NewService(repositoryService)

	// create an access token http handler and inject an access token service into this
	accessTokenHandler := http.NewHandler(accessTokenService)

	// Define a get token by id route
	router.GET("/oauth/access_token/:access_token_id", accessTokenHandler.GetById)

	// Run application on port 8082
	router.Run(":8082")
}
