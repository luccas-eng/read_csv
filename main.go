package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/read_csv/config"
	"github.com/read_csv/internal"
	"github.com/tinrab/retry"
)

func main() {

	args := os.Args[1:]
	fmt.Println(args)
	fileName := args[0]

	//load config with viper pkg
	config.LoadConfig()

	//connection string

	psqlconn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.C.GetString("db.user"), config.C.GetString("db.pwd"), config.C.GetString("db.host"), config.C.GetString("db.dbname"))

	fmt.Println(psqlconn)

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
	_, err := s.ProcessData(fileName)
	if err != nil {
		log.Println("s.ProcessData(): %w", err)
		panic(err)
	}
	elapsed1 := time.Since(start)

	log.Printf("took %.2f seconds to process data", elapsed1.Seconds())
	log.Println("process done, bye")

}
