package model

import "errors"

var (
	ErrOrderNotFound               = errors.New("order not found")
	ErrOrderCantBeCancelled        = errors.New("order paid or already cancelled")
	ErrItemNotFound                = errors.New("item not found")
	ErrOrderAlreadyPaidOrCancelled = errors.New("order already paid or cancelled")
)
