package middleware

import (
	"encoding/json"
	"errors"
	"github.com/LeonhardtDavid/xepelin-challenge/backend/internal/api"
	internalErrors "github.com/LeonhardtDavid/xepelin-challenge/backend/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandleErrors_NoError(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)

	HandleErrors(ctx)

	assert.Nil(t, ctx.Errors)
}

func Test_HandleErrors_BadRequestError(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)

	message := "Some bad request error"
	ctx.Error(&internalErrors.BadRequestError{
		Message: message,
	})

	HandleErrors(ctx)

	var response ErrorResponse
	err := json.NewDecoder(responseWriter.Body).Decode(&response)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, responseWriter.Code)
	assert.Equal(t, message, response.Error)
}

func Test_HandleErrors_UnauthorizedError(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)

	message := "Some unauthorized error"
	ctx.Error(&internalErrors.UnauthorizedError{
		Message: message,
	})

	HandleErrors(ctx)

	var response ErrorResponse
	err := json.NewDecoder(responseWriter.Body).Decode(&response)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, responseWriter.Code)
	assert.Equal(t, message, response.Error)
}

func Test_HandleErrors_InternalServerError(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)

	ctx.Error(errors.New("some unexpected error"))

	HandleErrors(ctx)

	var response ErrorResponse
	err := json.NewDecoder(responseWriter.Body).Decode(&response)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, responseWriter.Code)
	assert.Equal(t, "Something went wrong, sorry", response.Error)
}

func Test_RetrieveCustomer_ValidHeaderIsPresent(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)

	expectedCustomerId := uuid.New().String()
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set(api.CustomerIdHeader, expectedCustomerId)

	RetrieveCustomer(ctx)

	customerId, exists := ctx.Get(api.CustomerIdHeader)

	assert.True(t, exists)
	assert.Equal(t, expectedCustomerId, customerId.(uuid.UUID).String())
}

func Test_RetrieveCustomer_HeaderIsPresentWithInvalidFormat(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)

	customerId := "SomeInvalidId"
	ctx.Request = httptest.NewRequest(http.MethodPost, "/fail1", nil)
	ctx.Request.Header.Set(api.CustomerIdHeader, customerId)

	expectedError := internalErrors.UnauthorizedError{
		Message: "Unauthorized",
	}

	RetrieveCustomer(ctx)

	_, exists := ctx.Get(api.CustomerIdHeader)

	assert.False(t, exists)
	assert.Equal(t, &expectedError, ctx.Errors[0].Err)
}

func Test_RetrieveCustomer_HeaderNotIsPresent(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseWriter)

	ctx.Request = httptest.NewRequest(http.MethodPost, "/fail2", nil)

	expectedError := internalErrors.UnauthorizedError{
		Message: "Unauthorized",
	}

	RetrieveCustomer(ctx)

	_, exists := ctx.Get(api.CustomerIdHeader)

	assert.False(t, exists)
	assert.Equal(t, &expectedError, ctx.Errors[0].Err)
}
