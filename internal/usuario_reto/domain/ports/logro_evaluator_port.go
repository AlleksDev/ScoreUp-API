package ports

// LogroEvaluatorPort permite al módulo usuario_reto disparar la evaluación
// de logros tras completar un reto, sin acoplarse al módulo usuario_logro.
type LogroEvaluatorPort interface {
	EvaluateLogros(userID int64) ([]int64, error)
}
