package repository

import (
	"database/sql"
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/entities"
)

type UsuarioLogroMySQLRepository struct {
	conn *core.Conn_MySQL
}

func NewUsuarioLogroMySQLRepository(conn *core.Conn_MySQL) *UsuarioLogroMySQLRepository {
	return &UsuarioLogroMySQLRepository{conn: conn}
}

func (r *UsuarioLogroMySQLRepository) Save(ul *entities.UsuarioLogro) error {
	query := `
		INSERT INTO usuario_logros (id_usuario, id_logro)
		VALUES (?, ?)`

	_, err := r.conn.DB.Exec(query, ul.UserID, ul.LogroID)
	if err != nil {
		return fmt.Errorf("error al insertar usuario_logro: %w", err)
	}

	return nil
}

func (r *UsuarioLogroMySQLRepository) GetByUserID(userID int64) ([]*entities.UsuarioLogro, error) {
	query := `
		SELECT id_usuario, id_logro, fecha_obtenido
		FROM usuario_logros
		WHERE id_usuario = ?
		ORDER BY fecha_obtenido DESC`

	rows, err := r.conn.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo logros del usuario: %w", err)
	}
	defer rows.Close()

	var results []*entities.UsuarioLogro

	for rows.Next() {
		var ul entities.UsuarioLogro
		err := rows.Scan(
			&ul.UserID,
			&ul.LogroID,
			&ul.ObtainedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando usuario_logro: %w", err)
		}
		results = append(results, &ul)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando usuario_logros: %w", err)
	}

	return results, nil
}

func (r *UsuarioLogroMySQLRepository) Exists(userID int64, logroID int64) (bool, error) {
	query := `
		SELECT COUNT(1) FROM usuario_logros
		WHERE id_usuario = ? AND id_logro = ?`

	var count int
	err := r.conn.DB.QueryRow(query, userID, logroID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error verificando usuario_logro: %w", err)
	}

	return count > 0, nil
}

func (r *UsuarioLogroMySQLRepository) Delete(userID int64, logroID int64) error {
	query := `DELETE FROM usuario_logros WHERE id_usuario = ? AND id_logro = ?`

	_, err := r.conn.DB.Exec(query, userID, logroID)
	if err != nil {
		return fmt.Errorf("error eliminando usuario_logro: %w", err)
	}

	return nil
}
