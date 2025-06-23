package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"product-services/internal/logger"
	"product-services/internal/mocks"
	"product-services/internal/models"
	"product-services/internal/shared"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	env        = "prod"
	service    = "ProductService"
	ctxTimeOut = 5 * time.Second
)

func TestListCategories(t *testing.T) {
	createdAfter := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	const testLimit = 10

	mockRepo := new(mocks.MockCategoryRepository)
	mockUtil := new(mocks.MockSystemUtil)

	var logBuf bytes.Buffer
	logger := logger.NewLogger(env, service, &logBuf)
	h := NewCategoryHandler(mockRepo, mockUtil, logger, validator.New(), ctxTimeOut)

	t.Run("should respond with bad request if pagination validation fails", func(t *testing.T) {
		reqURL := "/categories?cursor=MjAyMy0wMS0wMVQwMDowMDowMFo&limit=ss"
		req := httptest.NewRequest(http.MethodGet, reqURL, strings.NewReader(""))
		rw := httptest.NewRecorder()

		h.ListCategories(rw, req)

		assert.Equal(t, http.StatusBadRequest, rw.Code)
		expectedResponse := `{
			"status":"error",
			"error": {
				"code": 1002,
				"message": "Invalid field format"
			}
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())

		mockRepo.AssertExpectations(t)
		mockUtil.AssertExpectations(t)

		// verify log content
		scanner := bufio.NewScanner(&logBuf)
		for scanner.Scan() {
			var entry map[string]interface{}
			err := json.Unmarshal(scanner.Bytes(), &entry)
			assert.NoError(t, err)
			assert.Equal(t, "error", entry["level"])
			assert.Equal(t, "ProductService", entry["service"])
			assert.Equal(t, "CategoryHandler.ListCategories", entry["op"])
			assert.Equal(t, float64(1002), entry["code"])
			assert.NotNil(t, entry["time"])
			errMsg := "strconv.ParseInt: parsing \"ss\": invalid syntax"
			assert.Equal(t, errMsg, entry["error"])
			assert.Equal(t, "Invalid field format", entry["message"])
			assert.Contains(t, entry["caller"], "internal/handlers/category_handler.go")
			assert.Nil(t, entry["details"])
		}
	})

	t.Run("should respond with internal server error if repo fails", func(t *testing.T) {
		dbError := errors.New("db query error")
		listOptions := shared.ListOptions{
			CreatedAfter: createdAfter,
			Limit:        testLimit,
		}
		mockRepo.On("ListCategories", mock.Anything, listOptions).
			Return(&models.ListCategoriesResult{}, dbError)

		reqURL := "/categories?cursor=MjAyMy0wMS0wMVQwMDowMDowMFo&limit=10"
		req := httptest.NewRequest(http.MethodGet, reqURL, strings.NewReader(""))
		rw := httptest.NewRecorder()

		h.ListCategories(rw, req)

		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		expectedResponse := `{
			"status":"error",
			"error": {
				"code": 1600,
				"message": "Internal server error"
			}
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())

		mockRepo.AssertExpectations(t)
		mockUtil.AssertExpectations(t)

		// verify log content
		scanner := bufio.NewScanner(&logBuf)
		for scanner.Scan() {
			var entry map[string]interface{}
			err := json.Unmarshal(scanner.Bytes(), &entry)
			assert.NoError(t, err)
			assert.Equal(t, "error", entry["level"])
			assert.Equal(t, "ProductService", entry["service"])
			assert.Equal(t, "CategoryHandler.ListCategories", entry["op"])
			assert.Equal(t, float64(1600), entry["code"])
			assert.NotNil(t, entry["time"])
			assert.Equal(t, "db query error", entry["error"])
			assert.Equal(t, "Internal server error", entry["message"])
			assert.Contains(t, entry["caller"], "internal/handlers/category_handler.go")
			assert.Equal(t, nil, entry["details"])
		}
	})
}
