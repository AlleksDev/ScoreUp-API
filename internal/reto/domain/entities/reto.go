package entities

import "time"

type Reto struct {
	ID            int64
	UserID        int64
	Subject       string
	Description   string
	Goal          int
	PointsAwarded int
	Deadline      *time.Time
	CreatedAt     time.Time
}
