package domain

import (
	"time"
)

type Lobby struct {
	ID        int64
	CreatedAt time.Time
	Location  Coordinate
}
