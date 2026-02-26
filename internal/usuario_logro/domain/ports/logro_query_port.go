package ports

import logroEntities "github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"

// LogroQueryPort permite al módulo usuario_logro consultar los logros disponibles
// sin acoplarse directamente al repositorio de logros.
type LogroQueryPort interface {
	GetAllLogros() ([]*logroEntities.Logro, error)
}
