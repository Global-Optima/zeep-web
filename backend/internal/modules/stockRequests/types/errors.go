package types

import "errors"

var (
	ErrExistingRequest      = errors.New("existing stock request found")
	ErrInsufficientStock    = errors.New("insufficient stock to fulfill the request")
	ErrOneRequestPerDay     = errors.New("only one request allowed per day")
	ErrStockRequestNotFound = errors.New("stock request not found")
)
