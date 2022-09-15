package product

import (
	"errors"
)

var (
	ErrUserNotExist        = errors.New("user does not exist")
	ErrProductNotFound     = errors.New("product is not found")
	ErrFileCannotBeCreated = errors.New("file cannot be created")
	ErrFileCannotBeRead    = errors.New("file cannot be read")
)
