package models

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Not found!")
	ErrConflict            = errors.New("Already exist!")
	ErrBadParamInput       = errors.New("Given Param is not valid")
)
