package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrSKUNotFound        = errors.New("sku not found")
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrAddItemToCart      = errors.New("add item to cart")
)
