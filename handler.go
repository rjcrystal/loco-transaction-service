package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func putTransaction(c *gin.Context) {
	var pt putTransactionRequest
	err := c.BindJSON(&pt) // Bind json request body to struct and validate request
	if err != nil {
		lg.Warn().Err(err).Msg("Failed to validate the request ")
		sendErrorResponse(errorResponse{Error: err.Error()}, c, 400)
		return
	}
	transactionID, _ := strconv.ParseInt(c.Param("transactionId"), 10, 64) // Get transaction ID from path

	lg.Debug().Interface("request", pt).Int64("transactionID", transactionID).Send()

	if transactionID == 0 {
		sendErrorResponse(errorResponse{Error: errMissingOrInvalidTransactionID.Error()}, c, 400)
		return
	}

	err = Store.saveTransaction(transactionID, pt.Amount, pt.Type, pt.Parent) // Save transaction in Postgres
	if err != nil {
		sendErrorResponse(errorResponse{Error: err.Error()}, c, 500)
		return
	}
	sendSuccessResponse(successResponse{Status: "OK"}, c)
}

func getTransactionSum(c *gin.Context) {
	transactionID, _ := strconv.ParseInt(c.Param("transactionId"), 10, 64)
	if transactionID == 0 {
		sendErrorResponse(errorResponse{Error: errMissingOrInvalidTransactionID.Error()}, c, 400)
		return
	}
	t, err := Store.getTransactionSum(transactionID)
	if err != nil {
		sendErrorResponse(errorResponse{Error: err.Error()}, c, 400)
		return
	}
	sendSuccessResponse(getTransactionSumResponse{Sum: t}, c)
}

func getTransactionsByType(c *gin.Context) {
	transactionType := c.Param("type")
	if len(transactionType) == 0 {
		sendErrorResponse(errorResponse{Error: errMissingOrInvalidTransactionID.Error()}, c, 400)
		return
	}
	transactions, err := Store.getTransactionByType(transactionType)
	if err != nil {
		sendErrorResponse(errorResponse{Error: err.Error()}, c, 400)
		return
	}
	if len(transactions) == 0 {
		sendSuccessResponse([]int64{}, c)
		return
	}
	sendSuccessResponse(transactions, c)
}

func getTransactionsDetails(c *gin.Context) {
	transactionID, _ := strconv.ParseInt(c.Param("transactionId"), 10, 64)
	if transactionID == 0 {
		sendErrorResponse(errorResponse{Error: errMissingOrInvalidTransactionID.Error()}, c, 400)
		return
	}
	t, err := Store.getTransactionByID(transactionID)
	if err != nil {
		sendErrorResponse(errorResponse{Error: err.Error()}, c, 400)
		return
	}
	sendSuccessResponse(t.toGetTransactionsDetailsResponse(), c)
}

// healthCheck : checks if DB is reachable
func healthCheck(c *gin.Context) {
	lg := GetLogger()
	err := Store.Ping()
	if err != nil {
		lg.Err(err).Msg("Failed to ping Postgres for healthcheck")
		sendErrorResponse(struct{}{}, c, 500)
		return
	}
	sendSuccessResponse(successResponse{Status: "OK"}, c)
}
