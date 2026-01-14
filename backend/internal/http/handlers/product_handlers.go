package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"math"

	"github.com/go-chi/chi/v5"

	"github.com/soydoradesu/product_discovery/internal/domain"
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

type searchResp struct {
	Items []domain.ProductSummary `json:"items"`
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	Total int64 `json:"total"`
	TotalPages int `json:"totalPages"`
}

func (h *ProductHandlers) Search(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()

	var params domain.SearchParams
	params.Q = strings.TrimSpace(qp.Get("q"))

	// category multi-value: category=1&category=2
	for _, s := range qp["category"] {
		id, err := strconv.ParseInt(s, 10, 64)
		if err == nil && id > 0 {
			params.CategoryID = append(params.CategoryID, id)
		}
	}

	if v := strings.TrimSpace(qp.Get("minPrice")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			f = math.Round(f*100) / 100
			params.MinPrice = &f
		}
	}
	if v := strings.TrimSpace(qp.Get("maxPrice")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			f = math.Round(f*100) / 100
			params.MaxPrice = &f
		}
	}

	if v := strings.TrimSpace(qp.Get("inStock")); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			params.InStock = &b
		}
	}

	params.Sort = qp.Get("sort")     
	params.Method = qp.Get("method")

	if v := qp.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.Page = n
		}
	}
	if v := qp.Get("pageSize"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			params.PageSize = n
		}
	}

	items, total, normalized, err := h.Products.Search(r.Context(), params)
	if err != nil {
		respond.Fail(w, http.StatusInternalServerError, "INTERNAL", "something went wrong")
		return
	}

	totalPages := int((total + int64(normalized.PageSize) - 1) / int64(normalized.PageSize))
	resp := searchResp{
		Items: items,
		Page: normalized.Page,
		PageSize: normalized.PageSize,
		Total: total,
		TotalPages: totalPages,
	}
	respond.JSON(w, http.StatusOK, resp)
}
