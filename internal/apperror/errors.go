package apperror

import (
	"errors"
)



var (
	ErrItemNotFound = errors.New("item not found")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)