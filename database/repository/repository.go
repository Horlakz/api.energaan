package repository

import "gorm.io/gorm"

type Respository struct{}

type Pageable struct {
	Page          int    `json:"page"`
	Size          int    `json:"size"`
	SortBy        string `json:"sortBy"`
	SortDirection string `json:"sortDir"`
	Search        string `json:"search"`
}

type Pagination struct {
	CurrentPage int64 `json:"currentPage"`
	TotalPages  int64 `json:"totalPages"`
	TotalItems  int64 `json:"totalItems"`
}

func GeneratePageable(database *gorm.DB) (pageable Pageable) {
	return pageable
}
