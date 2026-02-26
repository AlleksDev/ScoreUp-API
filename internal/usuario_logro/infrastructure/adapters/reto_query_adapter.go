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

func (a *RetoQueryAdapter) CountCompletedByUser(userID int64) (int, error) {
	query := `SELECT COUNT(1) FROM usuario_retos WHERE id_usuario = ? AND estado = 'completado'`

	var count int
	err := a.conn.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error contando retos completados: %w", err)
	}

	return count, nil
}
