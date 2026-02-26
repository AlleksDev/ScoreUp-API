package entities

import "time"

type UsuarioReto struct {
	UserID   int64
	RetoID   int64
	Progress int
	Status   string
	JoinedAt time.Time
}
