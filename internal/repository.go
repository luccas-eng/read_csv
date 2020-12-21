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
	InsertValues(data []string) error
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

func (m *databaseRepo) InsertValues(data []string) error {

	query := `insert into dataset_db.public.original_data (
				cpf,
				private,
				incomplete,
				last_purchase,
				avg_ticket,
				last_ticket,
				frequent_store,
				last_store
			 ) values ($1, $2, $3, $4, $5, $6, $7, $8);`

	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("m.db.Begin(): %w", err)
	}

	r, err := tx.Exec(query, data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7])
	if err != nil {
		return fmt.Errorf("tx.Exec(): %w", err)
	}

	if rows, err := r.RowsAffected(); rows == 0 {
		if err != nil {
			return fmt.Errorf("r.RowsAffected(): %w", err)
		}
		return fmt.Errorf("r.RowsAffected(): %d", rows)
	}

	tx.Commit()

	return nil
}
