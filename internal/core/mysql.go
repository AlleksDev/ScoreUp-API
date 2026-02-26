package core

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Conn_MySQL struct {
	DB *sql.DB
}

func GetMySQLPool() (*Conn_MySQL, error) {

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("la variable de entorno DB_DSN está vacía")
	}

	// Ejemplo DSN:
	// usuario:password@tcp(localhost:3306)/basedatos?parseTime=true

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error al verificar la conexión (ping): %w", err)
	}

	fmt.Println("Conexión a MySQL exitosa")
	return &Conn_MySQL{DB: db}, nil
}

// Wrapper simple para Exec (Insert, Update, Delete)
func (conn *Conn_MySQL) Execute(query string, values ...interface{}) (sql.Result, error) {
	result, err := conn.DB.Exec(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando query: %w", err)
	}
	return result, nil
}

func (conn *Conn_MySQL) Query(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := conn.DB.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error en select query: %w", err)
	}
	return rows, nil
}