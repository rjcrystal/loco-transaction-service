package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/gin-gonic/gin"
)

func SetUpRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	return router
}

var transactionReq = putTransactionRequest{
	Amount: 10,
	Type:   "ANIMALS",
}
var transactionID = rand.Int63n(10000)

var transactionAmountChild float64 = 100

var transactionReqWithParent = putTransactionRequest{
	Amount: transactionAmountChild,
	Type:   "ANIMALS",
	Parent: &transactionID,
}
var transactionIDChild = rand.Int63n(10000)

func TestMain(m *testing.M) {
	cleanupDB()
	m.Run()
	cleanupDB()
}

func cleanupDB() {
	// Setup code here
	initPostgresConnection()
	Store.deleteTransaction(transactionID, transactionReq.Type) // Delete the transaction if It Already exists
	// tear down later
}

// HAPPY FLOWS
func TestPutTransaction(t *testing.T) {
	r := SetUpRouter()
	r.PUT("/transactionservice/transaction/:transactionId", putTransaction)
	jsonValue, _ := json.Marshal(transactionReq)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTransactionByType(t *testing.T) {
	expectedResponse := fmt.Sprintf(`[%d]`, transactionID)

	r := SetUpRouter()
	r.GET("/transactionservice/types/:type", getTransactionsByType)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/transactionservice/types/%s", transactionReq.Type), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, expectedResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTransactionSumWithOutParent(t *testing.T) {
	expectedResponse := fmt.Sprintf(`{"sum":%g}`, transactionReq.Amount)

	r := SetUpRouter()
	r.GET("/transactionservice/sum/:transactionId", getTransactionSum)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/transactionservice/sum/%d", transactionID), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, expectedResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTransactionSaveWithParent(t *testing.T) {
	r := SetUpRouter()
	r.PUT("/transactionservice/transaction/:transactionId", putTransaction)
	jsonValue, _ := json.Marshal(transactionReqWithParent)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/%d", transactionIDChild), bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTransactionSumWithParents(t *testing.T) {
	expectedResponse := fmt.Sprintf(`{"sum":%g}`, transactionReq.Amount+transactionAmountChild)

	r := SetUpRouter()
	r.GET("/transactionservice/sum/:transactionId", getTransactionSum)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/transactionservice/sum/%d", transactionID), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, expectedResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTransactionDetails(t *testing.T) {
	expectedResponse := fmt.Sprintf(`{"amount":%g,"type":"%s","parent_id":null}`, transactionReq.Amount, transactionReq.Type)

	r := SetUpRouter()
	r.GET("/transactionservice/transaction/:transactionId", getTransactionsDetails)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/transactionservice/transaction/%d", transactionID), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, expectedResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHealthCheck(t *testing.T) {
	r := SetUpRouter()
	r.GET("/health", healthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// NEGATIVE FLOWS

func TestPutInvalidTransactionID(t *testing.T) {
	r := SetUpRouter()
	r.PUT("/transactionservice/transaction/:transactionId", putTransaction)
	jsonValue, _ := json.Marshal(transactionReq)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/hello"), bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPutInvalidTransaction(t *testing.T) {
	transactionReqInvalid := putTransactionRequest{
		Type: "ANIMALS",
	}

	r := SetUpRouter()
	r.PUT("/transactionservice/transaction/:transactionId", putTransaction)

	jsonValue, _ := json.Marshal(transactionReqInvalid)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/transactionservice/transaction/hello"), bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetTransactionByInvalidType(t *testing.T) {
	expectedResponse := `[]`

	r := SetUpRouter()
	r.GET("/transactionservice/types/:type", getTransactionsByType)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/transactionservice/types/%d", 11), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, expectedResponse, string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}
