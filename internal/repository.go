package internal

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" //psql
)

//Repository ...
type Repository interface {
	Close()
}

type databaseRepo struct {
	db *sql.DB
}

// NewRepository creates a new instance for database communication
func NewRepository(url string) (Repository, error) {

	// open db conn
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("sql.Open(): %w", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(10 * time.Second)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db.Ping(): %w", err)
	}

	return &databaseRepo{db}, nil
}

// Close closes the connection between the client and database
func (m *databaseRepo) Close() {
	m.db.Close()
}
