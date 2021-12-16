package domain

import "errors"

var ErrEventNotFound = errors.New("event not found")
var ErrDateBusy = errors.New("date is busy by other event")
