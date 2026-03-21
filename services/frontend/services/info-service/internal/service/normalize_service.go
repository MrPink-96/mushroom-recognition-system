package service

import "strings"

func normalizePagination(page, limit int) (int, int) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	return page, limit
}

func normalizeSortOrder(order string) string {
	if strings.ToLower(order) == "desc" {
		return "DESC"
	}
	return "ASC"
}
