package models

import (
	"time"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type Event struct {
	ID          ID
	UserID      ID
	Title       string
	Date        time.Time
	Duration    time.Duration
	Description string
}
