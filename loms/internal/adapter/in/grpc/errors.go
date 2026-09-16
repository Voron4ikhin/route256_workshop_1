package grpc

import (
	"errors"
	"route256/loms/internal/domain"

	"google.golang.org/grpc/codes"
)

var (
	ErrIncorrectUser         = errors.New("incorrect user")
	ErrIncorrectSKU          = errors.New("incorrect SKU")
	ErrIncorrectProductCount = errors.New("incorrect product count")
	ErrIncorrectQuantity     = errors.New("incorrect item quantity")
	ErrIncorrectOrderID      = errors.New("incorrect order id")
)

var businessErrorCode = []struct {
	err  error
	code codes.Code
}{
	{domain.ErrOrderNotFound, codes.NotFound},
	{domain.ErrStatusToPay, codes.FailedPrecondition},
	{domain.ErrReserveOrder, codes.FailedPrecondition},
}

func classifyBusinessError(err error) codes.Code {
	for _, m := range businessErrorCode {
		if errors.Is(err, m.err) {
			return m.code
		}
	}
	return codes.Internal
}
