package domain

import "time"

type Category struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
}

type ProductImage struct {
	URL string `json:"url"`
	Position int32 `json:"position"`
}

type Product struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Description string `json:"description"`
	Rating float64 `json:"rating"`
	InStock bool `json:"inStock"`
	CreatedAt time.Time `json:"createdAt"`

	Images []ProductImage `json:"images"`
	Categories []Category `json:"categories"`
}

type ProductSummary struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Rating float64 `json:"rating"`
	InStock bool `json:"inStock"`
	CreatedAt time.Time `json:"createdAt"`
	Thumbnail *string `json:"thumbnail,omitempty"`
	Categories []Category `json:"categories"`
}

type SearchParams struct {
	Q string
	CategoryID []int64
	MinPrice *int64
	MaxPrice *int64
	InStock *bool
	Sort string
	Method string
	Page int
	PageSize int
}