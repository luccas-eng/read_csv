package internal

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/read_csv/config"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func TestProcessData(t *testing.T) {

	start := time.Now()
	config.LoadTestConfig()

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.C.GetString("db.host"), config.C.GetInt64("db.port"), config.C.GetString("db.user"), config.C.GetString("db.pwd"), config.C.GetString("db.dbname"))

	var repository Repository
	retry.ForeverSleep(10*time.Second, func(_ int) (err error) {
		repository, err = NewRepository(psqlconn)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()

	service := NewService(repository)

	data, total, err := service.ProcessData(context.Background())
	if err != nil {
		log.Println(fmt.Errorf("service.ReadData(): %w", err))
	}

	elapsed := time.Since(start)
	log.Printf("%d lines processed in %.2f seconds", total, elapsed.Seconds())

	t.Log(data)

}

func TestSanitizeData(t *testing.T) {

	start := time.Now()
	config.LoadTestConfig()

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.C.GetString("db.host"), config.C.GetInt64("db.port"), config.C.GetString("db.user"), config.C.GetString("db.pwd"), config.C.GetString("db.dbname"))

	var repository Repository
	retry.ForeverSleep(10*time.Second, func(_ int) (err error) {
		repository, err = NewRepository(psqlconn)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()

	service := NewService(repository)

	ok, err := service.SanitizeData(context.Background())
	if err != nil {
		log.Println(fmt.Errorf("service.ReadData(): %w", err))
	}

	elapsed := time.Since(start)
	log.Printf("lines processed in %.2f seconds", elapsed.Seconds())
	log.Printf("process returned: %v", ok)

	t.Log(ok)

}
