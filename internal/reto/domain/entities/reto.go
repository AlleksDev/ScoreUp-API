package entities

import "time"

type Reto struct {
	ID            int64
	UserID        int64
	Subject       string
	Description   string
	Goal          int
	Progress      int
	PointsAwarded int
	Deadline      *time.Time
	Status        string
	CreatedAt     time.Time
}
