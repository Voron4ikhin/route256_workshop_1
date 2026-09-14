package handlers

import "errors"

var ErrIncorrectUser = errors.New("incorrect user")
var ErrIncorrectSKU = errors.New("incorrect SKU")
var ErrIncorrectProductCount = errors.New("incorrect product count")
var ErrIncorrectQuantity = errors.New("incorrect item quantity")
var ErrIncorrectOrderID = errors.New("incorrect order id")
