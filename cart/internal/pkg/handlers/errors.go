package handlers

import "errors"

var ErrIncorrectUser = errors.New("incorrect user")
var ErrIncorrectSKU = errors.New("incorrect SKU")
var ErrIncorrectQuantity = errors.New("incorrect item quantity")
