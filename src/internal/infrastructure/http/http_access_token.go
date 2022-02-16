// Package http handle http client operations
package http

import (
	"net/http"

	accessTokenDomain "github.com/esequielvirtuoso/book_store_oauth-api/src/domain/access_token"
	accesstoken "github.com/esequielvirtuoso/book_store_oauth-api/src/services/access_token"
	restErrors "github.com/esequielvirtuoso/go_utils_lib/rest_errors"
	"github.com/gin-gonic/gin"
)

func NewHandler(service accesstoken.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
}

type accessTokenHandler struct {
	service accesstoken.Service
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessTokenId := c.Param("access_token_id")

	accessToken, err := handler.service.GetById(accessTokenId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var request accessTokenDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := restErrors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	accessToken, err := handler.service.Create(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, accessToken)
}
