package repository

import "gorm.io/gorm"

type Respository struct{}

type Pageable struct {
	Page          int    `json:"page"`
	Size          int    `json:"size"`
	SortBy        string `json:"sort_by"`
	SortDirection string `json:"sort_dir"`
	Search        string `json:"search"`
}

type Pagination struct {
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
}

func GeneratePageable(database *gorm.DB) (pageable Pageable) {
	return pageable
}
