package handlers

import (
	"net/http"
	"time"

	"product-services/internal/interfaces"

	"github.com/go-playground/validator/v10"
)

type CategoryHandler struct {
	repo       interfaces.CategoryRepository
	util       interfaces.SystemUtil
	logger     interfaces.AppLogger
	validate   *validator.Validate
	ctxTimeOut time.Duration
}

func NewCategoryHandler(
	repo interfaces.CategoryRepository,
	util interfaces.SystemUtil,
	logger interfaces.AppLogger,
	validate *validator.Validate,
	ctxTimeOut time.Duration,
) *CategoryHandler {
	return &CategoryHandler{
		repo:       repo,
		util:       util,
		logger:     logger,
		validate:   validate,
		ctxTimeOut: ctxTimeOut,
	}
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandler.ListCategories"
	_, _, err := ParseAndValidatePagination(r)
	if err != nil {
		WriteErrorResponse(
			w,
			http.StatusBadRequest,
			ErrCodeInvalidFieldFormat,
			ErrMessageInvalidFieldFormat,
			err,
			nil,
			op,
			h.logger,
		)
		return
	}
}
