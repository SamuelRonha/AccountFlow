package domain

import "errors"

var (
	ErrAccountNotFound       = errors.New("account not found")
	ErrDocumentAlreadyUsed   = errors.New("document already in use")
	ErrOperationTypeNotFound = errors.New("operation type not found")
	ErrInvalidAmount         = errors.New("amount must not be zero")
	ErrInvalidField          = errors.New("invalid or missing required field")
)
