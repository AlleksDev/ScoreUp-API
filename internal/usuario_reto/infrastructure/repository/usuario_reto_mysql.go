package repository

import (
	"database/sql"
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/entities"
)

type UsuarioRetoMySQLRepository struct {
	conn *core.Conn_MySQL
}

func NewUsuarioRetoMySQLRepository(conn *core.Conn_MySQL) *UsuarioRetoMySQLRepository {
	return &UsuarioRetoMySQLRepository{conn: conn}
}

func (r *UsuarioRetoMySQLRepository) Save(ur *entities.UsuarioReto) error {
	query := `
		INSERT INTO usuario_retos (id_usuario, id_reto, progreso, estado)
		VALUES (?, ?, ?, ?)`

	_, err := r.conn.DB.Exec(query, ur.UserID, ur.RetoID, ur.Progress, ur.Status)
	if err != nil {
		return fmt.Errorf("error al insertar usuario_reto: %w", err)
	}

	return nil
}

func (r *UsuarioRetoMySQLRepository) GetByUserID(userID int64) ([]*entities.UsuarioReto, error) {
	query := `
		SELECT id_usuario, id_reto, progreso, estado, fecha_union
		FROM usuario_retos
		WHERE id_usuario = ?
		ORDER BY fecha_union DESC`

	rows, err := r.conn.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo retos del usuario: %w", err)
	}
	defer rows.Close()

	var results []*entities.UsuarioReto
	for rows.Next() {
		var ur entities.UsuarioReto
		err := rows.Scan(
			&ur.UserID,
			&ur.RetoID,
			&ur.Progress,
			&ur.Status,
			&ur.JoinedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando usuario_reto: %w", err)
		}
		results = append(results, &ur)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando retos del usuario: %w", err)
	}

	return results, nil
}

func (r *UsuarioRetoMySQLRepository) GetByRetoID(retoID int64) ([]*entities.UsuarioReto, error) {
	query := `
		SELECT id_usuario, id_reto, progreso, estado, fecha_union
		FROM usuario_retos
		WHERE id_reto = ?
		ORDER BY fecha_union DESC`

	rows, err := r.conn.DB.Query(query, retoID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo usuarios del reto: %w", err)
	}
	defer rows.Close()

	var results []*entities.UsuarioReto
	for rows.Next() {
		var ur entities.UsuarioReto
		err := rows.Scan(
			&ur.UserID,
			&ur.RetoID,
			&ur.Progress,
			&ur.Status,
			&ur.JoinedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando usuario_reto: %w", err)
		}
		results = append(results, &ur)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando usuarios del reto: %w", err)
	}

	return results, nil
}

func (r *UsuarioRetoMySQLRepository) GetByUserAndReto(userID int64, retoID int64) (*entities.UsuarioReto, error) {
	query := `
		SELECT id_usuario, id_reto, progreso, estado, fecha_union
		FROM usuario_retos
		WHERE id_usuario = ? AND id_reto = ?`

	var ur entities.UsuarioReto
	err := r.conn.DB.QueryRow(query, userID, retoID).Scan(
		&ur.UserID,
		&ur.RetoID,
		&ur.Progress,
		&ur.Status,
		&ur.JoinedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error obteniendo usuario_reto: %w", err)
	}

	return &ur, nil
}

func (r *UsuarioRetoMySQLRepository) Update(ur *entities.UsuarioReto) error {
	query := `
		UPDATE usuario_retos
		SET progreso = ?, estado = ?
		WHERE id_usuario = ? AND id_reto = ?`

	_, err := r.conn.DB.Exec(query, ur.Progress, ur.Status, ur.UserID, ur.RetoID)
	if err != nil {
		return fmt.Errorf("error al actualizar usuario_reto: %w", err)
	}

	return nil
}

func (r *UsuarioRetoMySQLRepository) Delete(userID int64, retoID int64) error {
	query := `DELETE FROM usuario_retos WHERE id_usuario = ? AND id_reto = ?`

	_, err := r.conn.DB.Exec(query, userID, retoID)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario_reto: %w", err)
	}

	return nil
}
