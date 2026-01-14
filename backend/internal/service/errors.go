package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")
)