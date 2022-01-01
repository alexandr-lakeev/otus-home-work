package domain

import "errors"

var (
	ErrEventNotFound    = errors.New("event not found")
	ErrDateBusy         = errors.New("date is busy by other event")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrPremissionDenied = errors.New("permission denied")
)
