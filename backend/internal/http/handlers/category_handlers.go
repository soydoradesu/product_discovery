package handlers

import (
	"net/http"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/http/respond"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type CategoryHandlers struct {
	Categories *service.CategoryService
}

type listCategoriesResp struct {
	Items []domain.Category `json:"items"`
}

func (h *CategoryHandlers) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.Categories.List(r.Context())
	if err != nil {
		respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "something went wrong")
		return
	}
	respond.JSON(w, http.StatusOK, listCategoriesResp{Items: items})
}