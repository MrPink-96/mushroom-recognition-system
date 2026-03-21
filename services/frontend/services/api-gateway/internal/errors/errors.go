package errors

import "errors"

var (
	ErrInternal        = errors.New("internal error")
	ErrMLUnavailable   = errors.New("ml service unavailable")
	ErrInfoUnavailable = errors.New("info service unavailable")
	ErrInvalidRequest  = errors.New("invalid request")
	ErrFileRequired    = errors.New("file required")
	ErrFileLarge       = errors.New("file too large")
	ErrInvalidFile     = errors.New("invalid file")
	ErrInvalidFileType = errors.New("only jpeg/png allowed")
	ErrReadFile        = errors.New("failed to read file")
)
