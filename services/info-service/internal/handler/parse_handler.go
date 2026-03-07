package handler

import (
	appErr "Info_Service/internal/errors"
	"strconv"
)

func parseID(idStr string) (int64, error) {
	if idStr == "" {
		return 0, appErr.ErrInvalidID
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, appErr.ErrInvalidID
	}

	if id <= 0 {
		return 0, appErr.ErrInvalidID
	}

	return id, nil
}

func parsePageAndLimit(pageStr, limitStr string) (int, int, error) {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, appErr.ErrInvalidPage
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, appErr.ErrInvalidLimit
	}

	return page, limit, nil
}
