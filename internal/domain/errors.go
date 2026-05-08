package domain

import "errors"

var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrDocumentAlreadyUsed   = errors.New("document already in use")
	ErrOperationTypeNotFound = errors.New("operation type not found")
	ErrInvalidAmount         = errors.New("amount must not be zero")
	ErrInvalidBalance        = errors.New("balance must not be negative")
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrInvalidField          = errors.New("invalid or missing required field")
	ErrInvalidFieldAmount    = errors.New("amount must be a valid number")
)
