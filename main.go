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

	_ = internal.NewService(repository)

	log.Println("service is running")

}
