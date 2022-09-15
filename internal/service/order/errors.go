package order

import (
	"errors"
)

var (
	ErrUserNotExist    = errors.New("user does not exist")
	ErrProductNotExist = errors.New("product does not exist")
)
