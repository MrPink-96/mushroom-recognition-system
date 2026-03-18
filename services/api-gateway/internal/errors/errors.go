package errors

import "errors"

var (
	ErrInternal        = errors.New("internal error")
	ErrMLUnavailable   = errors.New("ml service unavailable")
	ErrInfoUnavailable = errors.New("info service unavailable")
	ErrInvalidRequest  = errors.New("invalid request")
)
