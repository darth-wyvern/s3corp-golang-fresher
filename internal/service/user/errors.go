package user

import (
	"errors"
)

var (
	ErrEmailExisted           = errors.New("email existed")
	ErrUserIDExisted          = errors.New("user id existed")
	ErrUserNotFound           = errors.New("user is not found")
	ErrPasswordCannotBeHashed = errors.New("password cannot be hashed")
	ErrPasswordIncorrect      = errors.New("password is incorrect")
	ErrTokeCannotBeGenerated  = errors.New("token cannot be generated")
	ErrEmailNotExist          = errors.New("email does not exist")
)
