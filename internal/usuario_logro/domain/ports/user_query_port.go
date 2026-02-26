package ports

// UserQueryPort permite al módulo usuario_logro consultar datos del usuario
// sin acoplarse directamente al repositorio de usuarios.
type UserQueryPort interface {
	GetTotalScore(userID int64) (int, error)
}
