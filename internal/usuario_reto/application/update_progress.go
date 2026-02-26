package app

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/ports"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/repository"
)

type UpdateProgressResult struct {
	Completed     bool
	LogrosAwarded []int64
}

type UpdateProgress struct {
	repo           repository.UsuarioRetoRepository
	retoQuery      ports.RetoQueryPort
	userScore      ports.UserScorePort
	logroEvaluator ports.LogroEvaluatorPort
}

func NewUpdateProgress(
	repo repository.UsuarioRetoRepository,
	retoQuery ports.RetoQueryPort,
	userScore ports.UserScorePort,
	logroEvaluator ports.LogroEvaluatorPort,
) *UpdateProgress {
	return &UpdateProgress{
		repo:           repo,
		retoQuery:      retoQuery,
		userScore:      userScore,
		logroEvaluator: logroEvaluator,
	}
}

func (uc *UpdateProgress) Execute(userID int64, retoID int64, newProgress int) (*UpdateProgressResult, error) {
	// 1. Obtener el registro actual
	ur, err := uc.repo.GetByUserAndReto(userID, retoID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario_reto: %v", err)
	}
	if ur == nil {
		return nil, fmt.Errorf("el usuario no está unido a este reto")
	}
	if ur.Status == "completado" {
		return nil, fmt.Errorf("este reto ya fue completado")
	}

	// 2. Obtener la meta del reto
	goal, err := uc.retoQuery.GetGoal(retoID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener meta del reto: %v", err)
	}

	// 3. Actualizar progreso
	ur.Progress = newProgress

	result := &UpdateProgressResult{}

	// 4. Si alcanzó o superó la meta, marcar como completado
	if ur.Progress >= goal {
		ur.Status = "completado"
		result.Completed = true
	}

	if err := uc.repo.Update(ur); err != nil {
		return nil, fmt.Errorf("error al actualizar progreso: %v", err)
	}

	// 5. Si se completó, sumar puntos y evaluar logros
	if result.Completed {
		points, err := uc.retoQuery.GetPointsAwarded(retoID)
		if err != nil {
			return nil, fmt.Errorf("error al obtener puntos del reto: %v", err)
		}

		if err := uc.userScore.AddScore(userID, points); err != nil {
			return nil, fmt.Errorf("error al sumar puntos: %v", err)
		}

		logrosAwarded, err := uc.logroEvaluator.EvaluateLogros(userID)
		if err != nil {
			return nil, fmt.Errorf("error al evaluar logros: %v", err)
		}
		result.LogrosAwarded = logrosAwarded
	}

	return result, nil
}
