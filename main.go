package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/csv_reader/config"
	_ "github.com/lib/pq"
)

func main() {

	//load config with viper pkg
	config.LoadConfig()

	//connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.C.GetString("db.host"), config.C.GetInt64("db.port"), config.C.GetString("db.user"), config.C.GetString("db.pwd"), config.C.GetString("db.dbname"))

	// open database
	db, err := sql.Open("postgres", psqlconn)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(10 * time.Second)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("DB Connected!")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
