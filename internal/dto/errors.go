package dto

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInternalFailure = errors.New("internal failure")
	ErrBadRequest      = errors.New("bad request")
	ErrConflict        = errors.New("conflict")
)
