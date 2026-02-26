package entities

import "time"

type User struct {
	ID             int64
	Name         string
	Email          string
	Password       string
	TotalScore  int
	CreatedAt  time.Time
}