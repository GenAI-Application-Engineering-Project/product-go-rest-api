package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"product-services/internal/logger"
	"product-services/internal/mocks"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

const (
	env        = "prod"
	service    = "ProductService"
	ctxTimeOut = 5 * time.Second
)

func TestListCategories(t *testing.T) {
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
			assert.Contains(t, entry["caller"], "internal/handlers/category_handler.go:40")
			assert.Nil(t, entry["details"])
		}
	})
}
