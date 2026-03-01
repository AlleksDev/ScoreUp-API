package repository

import (
	"database/sql"
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"
)

type LogroMySQLRepository struct {
	conn *core.Conn_MySQL
}

func NewLogroMySQLRepository(conn *core.Conn_MySQL) *LogroMySQLRepository {
	return &LogroMySQLRepository{conn: conn}
}

func (r *LogroMySQLRepository) Save(logro *entities.Logro) error {
	query := `
		INSERT INTO logros (nombre, descripcion, puntos_requeridos, retos_requeridos)
		VALUES (?, ?, ?, ?)`

	result, err := r.conn.DB.Exec(
		query,
		logro.Name,
		logro.Description,
		logro.RequiredPoints,
		logro.RequiredRetos,
	)

	if err != nil {
		return fmt.Errorf("error al insertar logro: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error obteniendo ID insertado: %w", err)
	}

	logro.ID = id

	return nil
}

func (r *LogroMySQLRepository) GetByID(id int64) (*entities.Logro, error) {
	query := `
		SELECT id_logro, nombre, descripcion, puntos_requeridos, retos_requeridos
		FROM logros
		WHERE id_logro = ?`

	row := r.conn.DB.QueryRow(query, id)

	var logro entities.Logro

	err := row.Scan(
		&logro.ID,
		&logro.Name,
		&logro.Description,
		&logro.RequiredPoints,
		&logro.RequiredRetos,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando logro por ID: %w", err)
	}

	return &logro, nil
}

func (r *LogroMySQLRepository) GetAll() ([]*entities.Logro, error) {
	query := `
		SELECT id_logro, nombre, descripcion, puntos_requeridos, retos_requeridos
		FROM logros
		ORDER BY puntos_requeridos ASC`

	rows, err := r.conn.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo logros: %w", err)
	}
	defer rows.Close()

	var logros []*entities.Logro

	for rows.Next() {
		var logro entities.Logro

		err := rows.Scan(
			&logro.ID,
			&logro.Name,
			&logro.Description,
			&logro.RequiredPoints,
			&logro.RequiredRetos,
		)

		if err != nil {
			return nil, fmt.Errorf("error escaneando logro: %w", err)
		}

		logros = append(logros, &logro)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando logros: %w", err)
	}

	return logros, nil
}

func (r *LogroMySQLRepository) Update(logro *entities.Logro) error {
	query := `
		UPDATE logros
		SET nombre = ?, descripcion = ?, puntos_requeridos = ?, retos_requeridos = ?
		WHERE id_logro = ?`

	result, err := r.conn.DB.Exec(
		query,
		logro.Name,
		logro.Description,
		logro.RequiredPoints,
		logro.RequiredRetos,
		logro.ID,
	)

	if err != nil {
		return fmt.Errorf("error actualizando logro: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró logro con id %d", logro.ID)
	}

	return nil
}

func (r *LogroMySQLRepository) Delete(id int64) error {
	query := `DELETE FROM logros WHERE id_logro = ?`

	_, err := r.conn.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando logro: %w", err)
	}

	return nil
}
