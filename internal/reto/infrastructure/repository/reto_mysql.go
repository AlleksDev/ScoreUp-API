package repository

import (
	"database/sql"
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/entities"
)

type RetoMySQLRepository struct {
	conn *core.Conn_MySQL
}

func NewRetoMySQLRepository(conn *core.Conn_MySQL) *RetoMySQLRepository {
	return &RetoMySQLRepository{conn: conn}
}

func (r *RetoMySQLRepository) Save(reto *entities.Reto) error {
	query := `
		INSERT INTO retos (id_usuario, materia, descripcion, meta, puntos_otorgados, fecha_limite)
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.conn.DB.Exec(
		query,
		reto.UserID,
		reto.Subject,
		reto.Description,
		reto.Goal,
		reto.PointsAwarded,
		reto.Deadline,
	)

	if err != nil {
		return fmt.Errorf("error al insertar reto: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error obteniendo ID insertado: %w", err)
	}

	reto.ID = id

	return nil
}

func (r *RetoMySQLRepository) GetByID(id int64) (*entities.Reto, error) {
	query := `
		SELECT id_reto, id_usuario, materia, descripcion, meta, puntos_otorgados, fecha_limite, fecha_creacion
		FROM retos
		WHERE id_reto = ?`

	row := r.conn.DB.QueryRow(query, id)

	var reto entities.Reto
	var deadline sql.NullTime

	err := row.Scan(
		&reto.ID,
		&reto.UserID,
		&reto.Subject,
		&reto.Description,
		&reto.Goal,
		&reto.PointsAwarded,
		&deadline,
		&reto.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando reto por ID: %w", err)
	}

	if deadline.Valid {
		reto.Deadline = &deadline.Time
	}

	return &reto, nil
}

func (r *RetoMySQLRepository) GetAll() ([]*entities.Reto, error) {
	query := `
		SELECT id_reto, id_usuario, materia, descripcion, meta, puntos_otorgados, fecha_limite, fecha_creacion
		FROM retos
		ORDER BY fecha_creacion DESC`

	rows, err := r.conn.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo retos: %w", err)
	}
	defer rows.Close()

	var retos []*entities.Reto

	for rows.Next() {
		var reto entities.Reto
		var deadline sql.NullTime

		err := rows.Scan(
			&reto.ID,
			&reto.UserID,
			&reto.Subject,
			&reto.Description,
			&reto.Goal,
			&reto.PointsAwarded,
			&deadline,
			&reto.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error escaneando reto: %w", err)
		}

		if deadline.Valid {
			reto.Deadline = &deadline.Time
		}

		retos = append(retos, &reto)
	}

	return retos, nil
}

func (r *RetoMySQLRepository) GetByCreator(userID int64) ([]*entities.Reto, error) {
	query := `
		SELECT id_reto, id_usuario, materia, descripcion, meta, puntos_otorgados, fecha_limite, fecha_creacion
		FROM retos
		WHERE id_usuario = ?
		ORDER BY fecha_creacion DESC`

	rows, err := r.conn.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo retos del creador: %w", err)
	}
	defer rows.Close()

	var retos []*entities.Reto

	for rows.Next() {
		var reto entities.Reto
		var deadline sql.NullTime

		err := rows.Scan(
			&reto.ID,
			&reto.UserID,
			&reto.Subject,
			&reto.Description,
			&reto.Goal,
			&reto.PointsAwarded,
			&deadline,
			&reto.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error escaneando reto: %w", err)
		}

		if deadline.Valid {
			reto.Deadline = &deadline.Time
		}

		retos = append(retos, &reto)
	}

	return retos, nil
}

func (r *RetoMySQLRepository) Update(reto *entities.Reto) error {
	query := `
		UPDATE retos
		SET materia = ?, descripcion = ?, meta = ?, puntos_otorgados = ?, fecha_limite = ?
		WHERE id_reto = ?`

	result, err := r.conn.DB.Exec(
		query,
		reto.Subject,
		reto.Description,
		reto.Goal,
		reto.PointsAwarded,
		reto.Deadline,
		reto.ID,
	)

	if err != nil {
		return fmt.Errorf("error actualizando reto: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró reto con id %d", reto.ID)
	}

	return nil
}

func (r *RetoMySQLRepository) Delete(id int64) error {
	query := `DELETE FROM retos WHERE id_reto = ?`

	_, err := r.conn.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando reto: %w", err)
	}

	return nil
}
