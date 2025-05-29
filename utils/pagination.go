package utils

import (
	"math"
	"strconv"
)

type PaginationInfo struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

func CalculatePagination(total, page, limit int) PaginationInfo {
	if limit <= 0 {
		limit = 10 // Default
	}
	if page <= 0 {
		page = 1
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	// Ensure page doesn't exceed total pages
	if page > totalPages {
		page = totalPages
	}

	return PaginationInfo{
		CurrentPage:  page,
		TotalPages:   totalPages,
		TotalItems:   total,
		ItemsPerPage: limit,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}

func ParsePaginationParams(pageStr, limitStr string) (int, int) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Limit the maximum items per page
	validLimits := []int{10, 20, 30}
	isValidLimit := false
	for _, validLimit := range validLimits {
		if limit == validLimit {
			isValidLimit = true
			break
		}
	}

	if !isValidLimit && limitStr != "all" {
		limit = 10
	}

	if limitStr == "all" {
		limit = 0 // 0 means no limit
	}

	return page, limit
}