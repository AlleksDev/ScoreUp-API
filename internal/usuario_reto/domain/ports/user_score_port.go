package ports

// UserScorePort permite al módulo usuario_reto sumar puntos a un usuario
// cuando completa un reto, sin acoplarse directamente al módulo user.
type UserScorePort interface {
	AddScore(userID int64, points int) error
}
