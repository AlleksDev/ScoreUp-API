package adapters

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
)

// RetoQueryAdapter implementa ports.RetoQueryPort consultando directamente
// la tabla retos sin depender del módulo reto.
type RetoQueryAdapter struct {
	conn *core.Conn_MySQL
}

func NewRetoQueryAdapter(conn *core.Conn_MySQL) *RetoQueryAdapter {
	return &RetoQueryAdapter{conn: conn}
}

func (a *RetoQueryAdapter) GetPointsAwarded(retoID int64) (int, error) {
	query := `SELECT puntos_otorgados FROM retos WHERE id_reto = ?`

	var points int
	err := a.conn.DB.QueryRow(query, retoID).Scan(&points)
	if err != nil {
		return 0, fmt.Errorf("error obteniendo puntos del reto: %w", err)
	}

	return points, nil
}

func (a *RetoQueryAdapter) GetGoal(retoID int64) (int, error) {
	query := `SELECT meta FROM retos WHERE id_reto = ?`

	var goal int
	err := a.conn.DB.QueryRow(query, retoID).Scan(&goal)
	if err != nil {
		return 0, fmt.Errorf("error obteniendo meta del reto: %w", err)
	}

	return goal, nil
}
