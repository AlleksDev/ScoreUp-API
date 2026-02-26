package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/entities"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/ports"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/repository"
)

type EvaluateLogros struct {
	repo       repository.UsuarioLogroRepository
	userQuery  ports.UserQueryPort
	retoQuery  ports.RetoQueryPort
	logroQuery ports.LogroQueryPort
}

func NewEvaluateLogros(
	repo repository.UsuarioLogroRepository,
	userQuery ports.UserQueryPort,
	retoQuery ports.RetoQueryPort,
	logroQuery ports.LogroQueryPort,
) *EvaluateLogros {
	return &EvaluateLogros{
		repo:       repo,
		userQuery:  userQuery,
		retoQuery:  retoQuery,
		logroQuery: logroQuery,
	}
}

// Execute evalúa todos los logros disponibles para un usuario y otorga
// los que cumpla y aún no tenga asignados. Retorna los IDs de logros
// nuevos otorgados en esta evaluación.
func (uc *EvaluateLogros) Execute(userID int64) ([]int64, error) {
	// 1. Obtener puntos totales del usuario
	totalScore, err := uc.userQuery.GetTotalScore(userID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener puntos del usuario: %v", err)
	}

	// 2. Obtener retos completados del usuario
	completedRetos, err := uc.retoQuery.CountCompletedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("error al contar retos completados: %v", err)
	}

	// 3. Obtener todos los logros disponibles
	allLogros, err := uc.logroQuery.GetAllLogros()
	if err != nil {
		return nil, fmt.Errorf("error al obtener logros: %v", err)
	}

	// 4. Evaluar cada logro y otorgar si el usuario cumple los requisitos
	var awarded []int64
	for _, logro := range allLogros {
		// Verificar si ya tiene este logro
		exists, err := uc.repo.Exists(userID, logro.ID)
		if err != nil {
			return nil, fmt.Errorf("error al verificar logro existente: %v", err)
		}
		if exists {
			continue
		}

		// Verificar si cumple los requisitos
		meetsPoints := totalScore >= logro.RequiredPoints
		meetsRetos := completedRetos >= logro.RequiredRetos

		if meetsPoints && meetsRetos {
			ul := &entities.UsuarioLogro{
				UserID:  userID,
				LogroID: logro.ID,
			}
			if err := uc.repo.Save(ul); err != nil {
				return nil, fmt.Errorf("error al otorgar logro %d: %v", logro.ID, err)
			}
			awarded = append(awarded, logro.ID)
		}
	}

	return awarded, nil
}
