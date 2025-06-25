package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"product-services/internal/logger"
	"product-services/internal/mocks"

	"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
	const op = "TestHandler.TestMethod"
	t.Run("should write success response", func(t *testing.T) {
		data := map[string]string{
			"id":   "f2aa335f-6f91-4d4d-8057-53b0009bc376",
			"name": "test name",
		}

		var logBuf bytes.Buffer
		logger := logger.NewLogger(env, service, &logBuf)

		rw := httptest.NewRecorder()
		writeResponse(rw, http.StatusOK, op, data, logger)

		expectedResponse := `{
			"id": "f2aa335f-6f91-4d4d-8057-53b0009bc376",
			"name": "test name"
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())
		assert.Equal(t, "", logBuf.String())
	})

	t.Run("should respond with internal server error if encoding fails", func(t *testing.T) {
		type Node struct {
			Value string
			Next  *Node
		}
		data := &Node{Value: "A"}
		data.Next = data

		var logBuf bytes.Buffer
		logger := logger.NewLogger(env, service, &logBuf)

		rw := httptest.NewRecorder()
		writeResponse(rw, http.StatusOK, op, data, logger)

		expectedResponse := `{
			"status":"error",
			"error": {
				"message": "Internal Server Error"
			}
		}`
		assert.JSONEq(t, expectedResponse, rw.Body.String())
		// verify log content
		scanner := bufio.NewScanner(&logBuf)
		for scanner.Scan() {
			var entry map[string]interface{}
			err := json.Unmarshal(scanner.Bytes(), &entry)
			assert.NoError(t, err)
			assert.Equal(t, "error", entry["level"])
			assert.Equal(t, "ProductService", entry["service"])
			assert.Equal(t, op, entry["op"])
			assert.Equal(t, float64(1001), entry["code"])
			assert.NotNil(t, entry["time"])
			errMsg := "json: unsupported value: encountered a cycle via *handlers.Node"
			assert.Equal(t, errMsg, entry["error"])
			assert.Equal(t, "JSON encoding error", entry["message"])
			assert.Contains(t, entry["caller"], "internal/handlers/common_test.go")
			assert.Equal(t, nil, entry["details"])
		}
	})

	t.Run("should respond with internal server error if buffer writing fails", func(t *testing.T) {
		data := "test"
		dataBytes := []byte{0x22, 0x74, 0x65, 0x73, 0x74, 0x22, 0xa}
		err := errors.New("writer error")
		mockResponseWriter := new(mocks.MockHTTPResponseWriter)
		mockResponseWriter.On("Write", dataBytes).Return(0, err)
		mockResponseWriter.On("Header").Return(http.Header{})
		mockResponseWriter.On("WriteHeader", 200).Return()

		var logBuf bytes.Buffer
		logger := logger.NewLogger(env, service, &logBuf)

		writeResponse(mockResponseWriter, http.StatusOK, op, data, logger)
		// verify log content
		scanner := bufio.NewScanner(&logBuf)
		for scanner.Scan() {
			var entry map[string]interface{}
			err := json.Unmarshal(scanner.Bytes(), &entry)
			assert.NoError(t, err)
			assert.Equal(t, "error", entry["level"])
			assert.Equal(t, "ProductService", entry["service"])
			assert.Equal(t, op, entry["op"])
			assert.Equal(t, float64(1002), entry["code"])
			assert.NotNil(t, entry["time"])
			assert.Equal(t, "writer error", entry["error"])
			assert.Equal(t, "Failed response writer", entry["message"])
			assert.Contains(t, entry["caller"], "internal/handlers/common_test.go")
			assert.Equal(t, nil, entry["details"])
		}
	})
}
