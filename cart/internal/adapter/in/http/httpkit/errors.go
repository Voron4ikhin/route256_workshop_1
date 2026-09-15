package httpkit

import "errors"

var (
	ErrIncorrectUser     = errors.New("incorrect user")
	ErrIncorrectSKU      = errors.New("incorrect SKU")
	ErrIncorrectQuantity = errors.New("incorrect item quantity")
	ErrMethodNotAllowed  = errors.New("method not allowed")
)
