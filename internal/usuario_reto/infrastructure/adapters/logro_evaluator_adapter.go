package adapters

import (
	"fmt"

	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/application"
)

// LogroEvaluatorAdapter implementa ports.LogroEvaluatorPort invocando
// directamente el caso de uso EvaluateLogros del módulo usuario_logro.
type LogroEvaluatorAdapter struct {
	evaluateLogros *app.EvaluateLogros
}

func NewLogroEvaluatorAdapter(evaluateLogros *app.EvaluateLogros) *LogroEvaluatorAdapter {
	return &LogroEvaluatorAdapter{evaluateLogros: evaluateLogros}
}

func (a *LogroEvaluatorAdapter) EvaluateLogros(userID int64) ([]int64, error) {
	awarded, err := a.evaluateLogros.Execute(userID)
	if err != nil {
		return nil, fmt.Errorf("error evaluando logros: %w", err)
	}

	return awarded, nil
}
