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
	// Error codes
	ErrCodeValidationFailed     = 1000
	ErrCodeMissingRequiredField = 1001
	ErrCodeInvalidFieldFormat   = 1002
	ErrCodeInvalidRequestBody   = 1007
	ErrCodeRequestBodyTooLarge  = 1008
	ErrCodeInternalServerError  = 1600
	ErrCodeResourceNotFound     = 1300

	// Error code messages
	ErrMessageInvalidFieldFormat  = "Invalid field format"
	ErrMessageResourceNotFound    = "Resource not found"
	ErrMessageInternalServerError = "Internal server error"
	ErrMessageInvalidRequestBody  = "Invalid request body"
	ErrMessageRequestTooLarge     = "Request too large"
	ErrMessageValidationFailed    = "Validation failed"

	// Path params
	CursorParm = "cursor"
	LimitParam = "limit"

	StatusSuccess = "success"
	StatusError   = "error"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type HTTPErrorResponse struct {
	Status string `json:"status"`
	Error  Error  `json:"error"`
}

// DecodeCursorToTime decodes a base64 URL-safe string back into a time.Time
func DecodeCursorToTime(cursor string) (time.Time, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cursor encoding: %s", cursor)
	}

	t, err := time.Parse(time.RFC3339Nano, string(decodedBytes))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cursor time format: %s", cursor)
	}
	return t, nil
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
		return 0, nil
	}

	val, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}

func ParseAndValidatePagination(r *http.Request) (time.Time, int, error) {
	cursor, err := ParseCursor(r)
	if err != nil {
		return time.Time{}, 0, err
	}

	limit, err := ParseLimit(r)
	if err != nil {
		return time.Time{}, 0, err
	}

	return cursor, limit, nil
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
			// Ensure details is nil to avoid infinite recursion
			WriteErrorResponse(
				w,
				http.StatusInternalServerError,
				ErrCodeInternalServerError,
				ErrMessageInternalServerError,
				err,
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
			msg := "error writing response to client"
			appLogger.Err(err).Str("op", op).Int("code", ErrCodeInternalServerError).Msg(msg)
		}
	}
}

func WriteErrorResponse(
	w http.ResponseWriter,
	statusCode int,
	code int,
	message string,
	err error,
	details any,
	op string,
	logger interfaces.AppLogger,
) {
	appLogger := logger.Logger()
	appLogger.Err(err).Str("op", op).Int("code", code).Interface("details", details).Msg(message)
	resp := HTTPErrorResponse{
		Status: StatusError,
		Error: Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	}

	writeResponse(w, statusCode, op, resp, logger)
}
