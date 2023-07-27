package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type putTransactionRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Type   string  `json:"type" binding:"required"`
	Parent *int64  `json:"parent_id"`
}

type getTransactionsDetailsResponse struct {
	Amount float64 `json:"amount"`
	Type   string  `json:"type"`
	Parent *int64  `json:"parent_id"`
}

type transaction struct {
	ID     int64
	Amount float64
	Type   string
	Parent *int64
}

type getTransactionSumResponse struct {
	Sum float64 `json:"sum"`
}

func (t transaction) toGetTransactionsDetailsResponse() getTransactionsDetailsResponse {
	return getTransactionsDetailsResponse{
		Amount: t.Amount,
		Type:   t.Type,
		Parent: t.Parent,
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Status string `json:"status"`
}

var errMissingOrInvalidTransactionID = errors.New("Missing or Invalid Transaction ID in request path")

// sendSuccessResponse :
func sendSuccessResponse(payload interface{}, c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, payload)
}

// sendErrorResponse :
func sendErrorResponse(payload interface{}, c *gin.Context, code int) {
	c.AbortWithStatusJSON(code, payload)
}
