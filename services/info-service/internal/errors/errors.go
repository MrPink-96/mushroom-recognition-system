package errors

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal error")

	ErrInvalidPage       = errors.New("invalid page (must be positive integer)")
	ErrInvalidLimit      = errors.New("invalid limit (must be between 1 and 100)")
	ErrInvalidID         = errors.New("invalid id")
	ErrInvalidCategoryID = errors.New("invalid category id")
	ErrInvalidEdibility  = errors.New("invalid edibility (must be number)")
	ErrInvalidToxicity   = errors.New("invalid toxicity_max (must be number)")
)
