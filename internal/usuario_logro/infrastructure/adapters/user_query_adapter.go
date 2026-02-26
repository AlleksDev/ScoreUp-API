package adapters

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
)

// UserQueryAdapter implementa ports.UserQueryPort consultando directamente
// la tabla usuarios sin depender del módulo user.
type UserQueryAdapter struct {
	conn *core.Conn_MySQL
}

func NewUserQueryAdapter(conn *core.Conn_MySQL) *UserQueryAdapter {
	return &UserQueryAdapter{conn: conn}
}

func (a *UserQueryAdapter) GetTotalScore(userID int64) (int, error) {
	query := `SELECT puntos_totales FROM usuarios WHERE id_usuario = ?`

	var score int
	err := a.conn.DB.QueryRow(query, userID).Scan(&score)
	if err != nil {
		return 0, fmt.Errorf("error obteniendo puntos del usuario: %w", err)
	}

	return score, nil
}
