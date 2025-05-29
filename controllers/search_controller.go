package controllers

import (
	"encoding/json"
	"net/http"

	"goforum/models"
	"goforum/utils"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	tagName := r.URL.Query().Get("tag")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit := utils.ParsePaginationParams(pageStr, limitStr)

	// Build filters for search
	filters := models.ThreadFilters{
		Search:  query,
		TagName: tagName,
		Page:    page,
		Limit:   limit,
	}

	threads, total, err := models.GetThreads(filters)
	if err != nil {
		http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	pagination := utils.CalculatePagination(total, page, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"threads":    threads,
		"pagination": pagination,
		"query":      query,
		"tag":        tagName,
	})
}