package handlers

import (
	"context"
	"net/http"
	"time"

	"product-services/internal/interfaces"
	"product-services/internal/shared"

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
	createdAfter, limit, isValid := ParseAndValidatePagination(r, op, h.logger)
	if !isValid {
		WriteErrorResponse(
			w,
			http.StatusBadRequest,
			ErrMessageInvalidRequestParam,
			nil,
			op,
			h.logger,
		)
		return
	}

	listOptions := shared.ListOptions{
		CreatedAfter: createdAfter,
		Limit:        limit,
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.ctxTimeOut)
	defer cancel()

	result, err := h.repo.ListCategories(ctx, listOptions)
	if err != nil {
		WriteErrorResponse(
			w,
			http.StatusInternalServerError,
			ErrMessageInternalServerError,
			nil,
			op,
			h.logger,
		)
		return
	}

	WriteSuccessResponse(
		w,
		http.StatusOK,
		"Successfully fetched list of categories",
		result.Categories,
		&Pagination{
			HasMore:    result.HasMore,
			NextCursor: EncodeTimeToCursor(result.NextCursor),
		},
		op,
		h.logger,
	)
}
