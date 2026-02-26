package ports

// RetoQueryPort permite al módulo usuario_logro consultar datos de retos
// sin acoplarse directamente al repositorio de retos.
type RetoQueryPort interface {
	CountCompletedByUser(userID int64) (int, error)
}
