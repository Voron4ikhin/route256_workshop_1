package httpkit

import "errors"

var (
	ErrIncorrectUser         = errors.New("incorrect user")
	ErrIncorrectSKU          = errors.New("incorrect SKU")
	ErrIncorrectProductCount = errors.New("incorrect product count")
	ErrIncorrectQuantity     = errors.New("incorrect item quantity")
	ErrIncorrectOrderID      = errors.New("incorrect order id")
	ErrMethodNotAllowed      = errors.New("method not allowed")
)
