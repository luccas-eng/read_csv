package internal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/read_csv/internal/model"

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

func (m *databaseRepo) InsertValues(data model.Data) error {

	query := `insert into dataset_db.public.txtdata (
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

	r, err := tx.Exec(query, data.Cpf, data.Private, data.Incomplete, data.LastPurchase, data.AvgTicket, data.LastPurchase, data.FrequentStore, data.LastStore)
	if err != nil {
		return fmt.Errorf("tx.Exec(): %w", err)
	}

	if rows, err := r.RowsAffected(); rows == 0 {
		if err != nil {
			return fmt.Errorf("r.RowsAffected(): %w", err)
		}
		return fmt.Errorf("r.RowsAffected(): %d", rows)
	}

	return nil
}
