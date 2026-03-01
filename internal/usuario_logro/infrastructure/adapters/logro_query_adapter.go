package adapters

import (
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
	logroEntities "github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"
)

// LogroQueryAdapter implementa ports.LogroQueryPort consultando directamente
// la tabla logros sin depender del módulo logro.
type LogroQueryAdapter struct {
	conn *core.Conn_MySQL
}

func NewLogroQueryAdapter(conn *core.Conn_MySQL) *LogroQueryAdapter {
	return &LogroQueryAdapter{conn: conn}
}

func (a *LogroQueryAdapter) GetAllLogros() ([]*logroEntities.Logro, error) {
	query := `
		SELECT id_logro, nombre, descripcion, puntos_requeridos, retos_requeridos
		FROM logros`

	rows, err := a.conn.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo logros: %w", err)
	}
	defer rows.Close()

	var logros []*logroEntities.Logro

	for rows.Next() {
		var l logroEntities.Logro
		err := rows.Scan(
			&l.ID,
			&l.Name,
			&l.Description,
			&l.RequiredPoints,
			&l.RequiredRetos,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando logro: %w", err)
		}
		logros = append(logros, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando logros: %w", err)
	}

	return logros, nil
}
