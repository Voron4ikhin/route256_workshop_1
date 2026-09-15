package domain

import "errors"

var (
	ErrCreateOrder        = errors.New("cannot create order in OrderStorage")
	ErrOrderNotFound      = errors.New("cannot find order in OrderStorage")
	ErrReserveOrder       = errors.New("cannot reserve items in StockStorage")
	ErrReserveRemoveOrder = errors.New("cannot remove reserved items in StockStorage")
	ErrReserveCancelOrder = errors.New("cannot cancel reserve order")
	ErrStatusSetter       = errors.New("cannot set status in OrderStorage")
	ErrStatusToPay        = errors.New("cannot set status to pay because status is not ready for pay")
)
