package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/read_csv/config"
	"github.com/read_csv/internal"
	"github.com/tinrab/retry"
)

func main() {

	//load config with viper pkg
	config.LoadConfig()

	//connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.C.GetString("db.host"), config.C.GetInt64("db.port"), config.C.GetString("db.user"), config.C.GetString("db.pwd"), config.C.GetString("db.dbname"))

	var repository internal.Repository
	retry.ForeverSleep(10*time.Second, func(_ int) (err error) {
		repository, err = internal.NewRepository(psqlconn)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()

	s := internal.NewService(repository)

	log.Println("service is running")

	start := time.Now()
	log.Printf("started data processing at %s", start.Format("2006-01-02 03:04:05"))
	total, err := s.ProcessData()
	if err != nil {
		log.Println("s.ProcessData(): %w", err)
		panic(err)
	}
	elapsed1 := time.Since(start)

	start = time.Now()
	log.Printf("started data sanitizing at %s", start.Format("2006-01-02 03:04:05"))
	ok, err := s.SanitizeData()
	if err != nil {
		log.Println("s.SanitizeData(): %w", err)
		panic(err)
	}
	elapsed2 := time.Since(start)

	sanitized, err := s.CountSanitizedData()
	if err != nil {
		log.Println("s.CountSanitizedData(): %w", err)
		panic(err)
	}

	reliability := (float64(sanitized) / float64(total)) * float64(100)

	log.Printf("took %.2f seconds to process data with %d lines processed", elapsed1.Seconds(), total)
	log.Printf("took %.2f seconds to sanitize data - %v return", elapsed2.Seconds(), ok)
	log.Printf("reliability tax %.2f of sanitized ones", reliability)
	log.Println("process done, bye")

}
