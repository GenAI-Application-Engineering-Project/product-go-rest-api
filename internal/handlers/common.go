package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"product-services/internal/interfaces"
)

const (
	// Defaults
	DefaultLimit = 20

	// Error codes
	ErrCodeInvalidRequestParam  = 1000
	ErrCodeJSONEncoding         = 1001
	ErrCodeFailedResponseWriter = 1002

	// Error code messages
	ErrMessageInvalidRequestParam  = "Invalid request param"
	ErrMessageJSONEncoding         = "JSON encoding error"
	ErrMessageFailedResponseWriter = "Failed response writer"

	// http error Messages
	ErrMessageInternalServerError = "Internal Server Error"
	ErrMessageBadRequest          = "Bad Request"

	// Path params
	CursorParm = "cursor"
	LimitParam = "limit"

	StatusSuccess = "success"
	StatusError   = "error"
)

type Error struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type HTTPErrorResponse struct {
	Status string `json:"status"`
	Error  Error  `json:"error"`
}

type Pagination struct {
	HasMore    bool   `json:"has_more,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type HTTPSuccessResponse struct {
	Status     string      `json:"status"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Message    string      `json:"message"`
}

// DecodeCursorToTime decodes a base64 URL-safe string back into a time.Time
func DecodeCursorToTime(cursor string) (time.Time, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cursor encoding: `%s`, error: %v", cursor, err)
	}

	t, err := time.Parse(time.RFC3339Nano, string(decodedBytes))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cursor time format: `%s`, error: %v", cursor, err)
	}
	return t, nil
}

// EncodeTimeToCursor encodes a time.Time into a base64 URL-safe string
func EncodeTimeToCursor(t time.Time) string {
	timeStr := t.UTC().Format(time.RFC3339Nano)
	return base64.RawURLEncoding.EncodeToString([]byte(timeStr))
}

func ParseCursor(r *http.Request) (time.Time, error) {
	cursorStr := r.URL.Query().Get(CursorParm)
	if cursorStr == "" {
		return time.Time{}, nil
	}

	createdAfter, err := DecodeCursorToTime(cursorStr)
	if err != nil {
		return time.Time{}, err
	}
	return createdAfter, nil
}

func ParseLimit(r *http.Request) (int, error) {
	limitStr := r.URL.Query().Get(LimitParam)
	if limitStr == "" {
		return DefaultLimit, nil
	}

	val, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid limit value: `%s`, error: %v", limitStr, err)
	}

	return int(val), nil
}

func ParseAndValidatePagination(
	r *http.Request,
	op string,
	logger interfaces.AppLogger,
) (time.Time, int, bool) {
	cursor, err := ParseCursor(r)
	if err != nil {
		appLogger := logger.Logger()
		appLogger.Err(err).
			Str("op", op).
			Int("code", ErrCodeInvalidRequestParam).
			Msg(ErrMessageInvalidRequestParam)
		return time.Time{}, 0, false
	}

	limit, err := ParseLimit(r)
	if err != nil {
		appLogger := logger.Logger()
		appLogger.Err(err).
			Str("op", op).
			Int("code", ErrCodeInvalidRequestParam).
			Msg(ErrMessageInvalidRequestParam)
		return time.Time{}, 0, false
	}

	return cursor, limit, true
}

func writeResponse(
	w http.ResponseWriter,
	statusCode int,
	op string,
	details any,
	logger interfaces.AppLogger,
) {
	var buf bytes.Buffer
	if details != nil {
		err := json.NewEncoder(&buf).Encode(details)
		if err != nil {
			appLogger := logger.Logger()
			appLogger.Err(err).
				Str("op", op).
				Int("code", ErrCodeJSONEncoding).
				Msg(ErrMessageJSONEncoding)
			WriteErrorResponse(
				w,
				http.StatusInternalServerError,
				ErrMessageInternalServerError,
				nil,
				op,
				logger,
			)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Write response body
	if buf.Len() > 0 {
		if _, err := buf.WriteTo(w); err != nil {
			appLogger := logger.Logger()
			appLogger.Err(err).
				Str("op", op).
				Int("code", ErrCodeFailedResponseWriter).
				Msg(ErrMessageFailedResponseWriter)
		}
	}
}

func WriteErrorResponse(
	w http.ResponseWriter,
	statusCode int,
	message string,
	details any,
	op string,
	logger interfaces.AppLogger,
) {
	resp := HTTPErrorResponse{
		Status: StatusError,
		Error: Error{
			Message: message,
			Details: details,
		},
	}

	writeResponse(w, statusCode, op, resp, logger)
}

func WriteSuccessResponse(
	w http.ResponseWriter,
	statusCode int,
	message string,
	data any,
	pagination *Pagination,
	op string,
	logger interfaces.AppLogger,
) {
	resp := HTTPSuccessResponse{
		Status:     StatusSuccess,
		Data:       data,
		Pagination: pagination,
		Message:    message,
	}

	writeResponse(w, statusCode, op, resp, logger)
}
