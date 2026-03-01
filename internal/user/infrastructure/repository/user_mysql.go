package repository

import (
	"database/sql"
	"fmt"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
	"github.com/AlleksDev/ScoreUp-API/internal/user/domain/entities"
)

type UserMySQLRepository struct {
	conn *core.Conn_MySQL
}

func NewUserMySQLRepository(conn *core.Conn_MySQL) *UserMySQLRepository {
	return &UserMySQLRepository{conn: conn}
}

func (r *UserMySQLRepository) Save(user *entities.User) error {
	query := `
		INSERT INTO usuarios (nombre, email, password)
		VALUES (?, ?, ?)`

	result, err := r.conn.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
	)

	if err != nil {
		return fmt.Errorf("error al insertar usuario: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error obteniendo ID insertado: %w", err)
	}

	user.ID = id

	return nil
}

func (r *UserMySQLRepository) GetByEmail(email string) (*entities.User, error) {
	query := `
		SELECT id_usuario, nombre, email, password, puntos_totales, fecha_registro
		FROM usuarios
		WHERE email = ?`

	row := r.conn.DB.QueryRow(query, email)

	var user entities.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.TotalScore,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando usuario por email: %w", err)
	}

	return &user, nil
}

func (r *UserMySQLRepository) GetByID(id int64) (*entities.User, error) {
	query := `
		SELECT id_usuario, nombre, email, password, puntos_totales, fecha_registro
		FROM usuarios
		WHERE id_usuario = ?`

	row := r.conn.DB.QueryRow(query, id)

	var user entities.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.TotalScore,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando usuario por ID: %w", err)
	}

	return &user, nil
}

func (r *UserMySQLRepository) Update(user *entities.User) error {
	query := `
		UPDATE usuarios 
		SET nombre = ?, email = ?, password = ?, puntos_totales = ?
		WHERE id_usuario = ?`

	result, err := r.conn.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.TotalScore,

		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error actualizando usuario: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró usuario con id %d", user.ID)
	}

	return nil
}

func (r *UserMySQLRepository) Delete(id int64) error {
	query := `DELETE FROM usuarios WHERE id_usuario = ?`

	_, err := r.conn.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando usuario: %w", err)
	}

	return nil
}

func (r *UserMySQLRepository) GetRank() ([]entities.RankUser, error) {
	query := `
		SELECT id_usuario, nombre, puntos_totales
		FROM usuarios
		ORDER BY puntos_totales DESC
		LIMIT 10`

	rows, err := r.conn.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ranking: %w", err)
	}
	defer rows.Close()

	var users []entities.RankUser

	for rows.Next() {
		var user entities.RankUser

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.TotalScore,
		)

		if err != nil {
			return nil, fmt.Errorf("error escaneando ranking: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando ranking: %w", err)
	}

	return users, nil
}
