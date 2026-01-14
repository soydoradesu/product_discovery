package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/soydoradesu/product_discovery/internal/http/respond"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type ProductHandlers struct {
	Products *service.ProductService
}

func (h *ProductHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		respond.Fail(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid product id")
		return
	}

	p, err := h.Products.GetByID(r.Context(), id)
	if err != nil {
		if err == service.ErrProductNotFound {
			respond.Fail(w, http.StatusNotFound, "NOT_FOUND", "product not found")
			return
		}
		respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "something went wrong")
		return
	}

	respond.JSON(w, http.StatusOK, p)
}
