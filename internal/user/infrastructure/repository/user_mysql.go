package repository

import (
	"database/sql"
	"fmt"
	"strings"

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
		INSERT INTO users (username, name, email, password, phone, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.conn.DB.Exec(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error al insertar usuario: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error obteniendo ID insertado: %v", err)
	}

	user.ID = id
	return nil
}

func (r *UserMySQLRepository) GetByEmail(email string) (*entities.User, error) {
	query := `SELECT id, username, name, email, password, phone, created_at, updated_at FROM users WHERE email = ?`

	row := r.conn.DB.QueryRow(query, email)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando usuario por email: %v", err)
	}

	return &user, nil
}

func (r *UserMySQLRepository) GetByID(id int64) (*entities.User, error) {
	query := `SELECT id, username, name, email, password, phone, created_at, updated_at FROM users WHERE id = ?`

	row := r.conn.DB.QueryRow(query, id)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error buscando usuario por ID: %v", err)
	}

	return &user, nil
}

func (r *UserMySQLRepository) Update(user *entities.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, name = ?, phone = ?, password = ?, updated_at = ? 
		WHERE id = ?`

	result, err := r.conn.DB.Exec(
		query,
		user.Email,
		user.Name,
		user.Password,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error actualizando usuario: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró usuario con id %d para actualizar", user.ID)
	}

	return nil
}

func (r *UserMySQLRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.conn.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando usuario: %v", err)
	}
	return nil
}

func (r *UserMySQLRepository) GetUsersByIDs(ids []int64) ([]entities.User, error) {

	if len(ids) == 0 {
		return []entities.User{}, nil
	}

	// Construir (?, ?, ?, ?)
	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
        SELECT id, username, name, email, phone 
        FROM users 
        WHERE id IN (%s)`, placeholders)

	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}

	rows, err := r.conn.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error buscando usuarios por lote: %v", err)
	}
	defer rows.Close()

	var users []entities.User

	for rows.Next() {
		var u entities.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("error escaneando usuario: %v", err)
		}
		users = append(users, u)
	}

	return users, nil
}
