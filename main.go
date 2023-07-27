package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var lg zerolog.Logger

func init() {
	lg = GetLogger()

	err := godotenv.Load()
	if err != nil {
		lg.Fatal().Msg("Error loading .env file")
	}
	initPostgresConnection()
}

func main() {
	gin.DisableConsoleColor()
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.PUT("/transactionservice/transaction/:transactionId", putTransaction)
	router.GET("/transactionservice/types/:type", getTransactionsByType)
	router.GET("/transactionservice/sum/:transactionId", getTransactionSum)
	router.GET("/transactionservice/transaction/:transactionId", getTransactionsDetails)
	router.GET("/health", healthCheck)

	router.Run(":80")
}
