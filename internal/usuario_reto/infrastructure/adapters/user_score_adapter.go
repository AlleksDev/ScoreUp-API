package adapters

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
)

// UserScoreAdapter implementa ports.UserScorePort sumando puntos directamente
// en la tabla usuarios sin depender del módulo user.
type UserScoreAdapter struct {
	conn *core.Conn_MySQL
}

func NewUserScoreAdapter(conn *core.Conn_MySQL) *UserScoreAdapter {
	return &UserScoreAdapter{conn: conn}
}

func (a *UserScoreAdapter) AddScore(userID int64, points int) error {
	query := `UPDATE usuarios SET puntos_totales = puntos_totales + ? WHERE id_usuario = ?`

	_, err := a.conn.DB.Exec(query, points, userID)
	if err != nil {
		return fmt.Errorf("error al sumar puntos al usuario: %w", err)
	}

	return nil
}
