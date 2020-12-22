package internal

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/spf13/cast"

	_ "github.com/lib/pq" //psql
	"github.com/read_csv/internal/model"
	"github.com/read_csv/internal/util"
)

//Repository ...
type Repository interface {
	Close()
	InsertValues(data []string) error
	InsertSanitizedData(data []interface{}) error
	GetTotalLines() (totalLines int, err error)
	GetData(limit, offset int) (data []*model.Data, err error)
	CountSanitizedData() (totalLines int, err error)
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
	// log.Println("data processed")
	return nil
}

func (m *databaseRepo) InsertSanitizedData(data []interface{}) error {

	query := `insert into dataset_db.public.copy_data (
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

	var (
		cpf, lastPurchase, frequentStore, lastStore util.NullString
		avgTicket, lastTicket                       util.NullFloat64
	)

	a, er := cast.ToStringE(data[0])
	if er != nil {
		return fmt.Errorf("cast.ToStringE(data[0])")
	}
	b, er := cast.ToStringE(data[1])
	if er != nil {
		return fmt.Errorf("cast.ToStringE(data[1])")
	}
	c, er := cast.ToStringE(data[2])
	if er != nil {
		return fmt.Errorf("cast.ToStringE(data[2])")
	}
	d, er := cast.ToStringE(data[3])
	if er != nil {
		return fmt.Errorf("cast.ToStringE(data[3])")
	}
	e, er := cast.ToFloat64E(data[4])
	if er != nil {
		return fmt.Errorf("cast.ToStringE(data[4])")
	}
	f, er := cast.ToFloat64E(data[5])
	if er != nil {
		return fmt.Errorf("cast.ToStringE(data[5])")
	}
	g, er := cast.ToStringE(data[6])
	if er != nil {
		return fmt.Errorf("cast.ToStringE(data[6])")
	}
	h, er := cast.ToStringE(data[7])
	if er != nil {
		return fmt.Errorf("cast.ToStringE(data[7])")
	}

	r, err := tx.Exec(query, cpf.NewString(a), b, c, lastPurchase.NewString(d), avgTicket.NewFloat64(e), lastTicket.NewFloat64(f), frequentStore.NewString(g), lastStore.NewString(h))
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
	// log.Println("data processed")
	return nil
}

func (m *databaseRepo) GetTotalLines() (totalLines int, err error) {
	query := `select count(*) from dataset_db.public.original_data;`
	row := m.db.QueryRow(query)
	if err := row.Scan(&totalLines); err != nil {
		return 0, fmt.Errorf("row.Scan(): %w", err)
	}
	return
}

func (m *databaseRepo) CountSanitizedData() (totalLines int, err error) {
	query := `select count(*) from dataset_db.public.copy_data;`
	row := m.db.QueryRow(query)
	if err := row.Scan(&totalLines); err != nil {
		return 0, fmt.Errorf("row.Scan(): %w", err)
	}
	return
}

func (m *databaseRepo) GetData(limit, offset int) (data []*model.Data, err error) {
	query := `select cpf, private, incomplete, last_purchase, avg_ticket, last_ticket, frequent_store, last_store from dataset_db.public.original_data limit $1 offset $2;`

	rows, e := m.db.Query(query, limit, offset)
	if e != nil {
		err = fmt.Errorf("m.db.Query(): %w", e)
		return
	}
	defer rows.Close()

	for rows.Next() {
		dd := &model.Data{}

		e := rows.Scan(&dd.Cpf, &dd.Private, &dd.Incomplete, &dd.LastPurchase, &dd.AvgTicket, &dd.LastTicket, &dd.FrequentStore, &dd.LastStore)

		if e != nil {
			if e == sql.ErrNoRows {
				err = nil
				return
			}
			err = fmt.Errorf("rows.Scan(): %w", e)
			return
		}

		data = append(data, dd)
	}

	return
}
