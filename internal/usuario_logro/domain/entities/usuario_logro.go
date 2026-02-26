package entities

import "time"

type UsuarioLogro struct {
	UserID     int64
	LogroID    int64
	ObtainedAt time.Time
}
