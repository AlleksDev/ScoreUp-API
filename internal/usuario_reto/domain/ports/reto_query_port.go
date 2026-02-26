package ports

// RetoQueryPort permite al módulo usuario_reto consultar datos del reto
// (puntos otorgados, meta) sin acoplarse directamente al módulo reto.
type RetoQueryPort interface {
	GetPointsAwarded(retoID int64) (int, error)
	GetGoal(retoID int64) (int, error)
}
