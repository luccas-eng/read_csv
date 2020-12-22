package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/miguelpragier/handy"
	"golang.org/x/sync/semaphore"

	"github.com/read_csv/internal/model"
)

//Service interface implement services running between clients
type Service interface {
	ProcessData() (total int, err error)
	SanitizeData() (ok bool, err error)
	CountSanitizedData() (totalLines int, err error)
}

//InternalService struct implements repo
type InternalService struct {
	repository Repository
}

//NewService implements the repo service
func NewService(r Repository) Service {
	return &InternalService{r}
}

//ProcessData ...
func (s *InternalService) ProcessData() (total int, err error) {

	file, err := os.Open("./external/base_teste.txt")
	if err != nil {
		return 0, fmt.Errorf("os.Open(): %w", err)
	}

	defer file.Close()

	var counter int
	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)
	for {

		var (
			buffer   bytes.Buffer
			l        []byte
			isPrefix bool
		)

		for {

			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				break
			}

			// If we're at the EOF, break.
			if err != nil {
				if err != io.EOF {
					return 0, fmt.Errorf("%w", err)
				}
				break
			}
		}

		if err == io.EOF {
			err = nil
			break
		}

		//Discard first line
		if ok, err := regexp.MatchString("CPF", buffer.String()); ok && err == nil {
			continue
		}

		r := regexp.MustCompile("[^\\s]+")

		v := r.FindAllString(buffer.String(), -1)

		err = s.repository.InsertValues(v)
		if err != nil {
			return 0, fmt.Errorf("s.repository.InsertValues(): %w", err)
		}

		counter++

	}

	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("process failed with error: %w", err)
	}

	total = counter

	return
}

//SanitizeData ...
func (s *InternalService) SanitizeData() (ok bool, err error) {

	var (
		pages         int
		limit, offset = 500, 0
	)

	total, err := s.repository.GetTotalLines() //validate total of lines inserted
	if err != nil {
		err = fmt.Errorf("s.repository.GetTotalLines(): %w", err)
		return
	}

	pages = total / limit                    //get total pages
	if modDiv := total % limit; modDiv > 0 { //search more pages
		pages++
	}

	for page := 1; page <= pages; page++ {

		if page > 1 {
			offset = (limit * page) - limit
		}

		sanitizeData := func(values []*model.Data, limit, offset int) (totalErrors int, err error) {

			wg := &sync.WaitGroup{}
			wg.Add(len(values))

			for i, v := range values {

				var data []interface{}

				if strings.ToLower(strings.TrimSpace(v.Cpf)) != "null" {
					if handy.CheckCPF(v.Cpf) {
						v.Cpf = handy.OnlyDigits(v.Cpf)
						data = append(data, v.Cpf)
					}
				} else {
					data = append(data, "")
				}

				if strings.ToLower(strings.TrimSpace(v.Private)) != "null" {
					v.Private = handy.OnlyLettersAndNumbers(v.Private)
					data = append(data, v.Private)
				} else {
					data = append(data, "")
				}

				if strings.ToLower(strings.TrimSpace(v.Incomplete)) != "null" {
					v.Incomplete = handy.OnlyLettersAndNumbers(v.Incomplete)
					data = append(data, v.Incomplete)
				} else {
					data = append(data, "")
				}

				if strings.ToLower(strings.TrimSpace(v.LastPurchase)) != "null" {
					v.LastPurchase = handy.OnlyLettersAndNumbers(v.LastPurchase)
					data = append(data, v.LastPurchase)
				} else {
					data = append(data, "")
				}

				var avgTicket float64
				if strings.ToLower(strings.TrimSpace(v.AvgTicket)) != "null" {
					v.AvgTicket = strings.ReplaceAll(v.AvgTicket, ",", ".")
					avgTicket, err = strconv.ParseFloat(v.AvgTicket, 64)
					if err != nil {
						err = fmt.Errorf("strconv.ParseFloat(): %w", err)
						return
					}
					data = append(data, avgTicket)
				} else {
					data = append(data, avgTicket)
				}

				var lastTicket float64
				if strings.ToLower(strings.TrimSpace(v.LastTicket)) != "null" {
					v.LastTicket = strings.ReplaceAll(v.LastTicket, ",", ".")
					lastTicket, err = strconv.ParseFloat(v.LastTicket, 64)
					if err != nil {
						err = fmt.Errorf("strconv.ParseFloat(): %w", err)
						return
					}
					data = append(data, lastTicket)
				} else {
					data = append(data, lastTicket)
				}

				if strings.ToLower(strings.TrimSpace(v.FrequentStore)) != "null" {
					if handy.CheckCNPJ(v.FrequentStore) {
						v.FrequentStore = handy.OnlyDigits(v.FrequentStore)
					}
					data = append(data, v.FrequentStore)
				} else {
					data = append(data, "")
				}

				if strings.ToLower(strings.TrimSpace(v.LastStore)) != "null" {
					if handy.CheckCNPJ(v.LastStore) {
						v.LastStore = handy.OnlyDigits(v.LastStore)
					}
					data = append(data, v.LastStore)
				} else {
					data = append(data, "")
				}

				var cErr = make(chan error)

				go func() {
					e := s.repository.InsertSanitizedData(data)
					defer wg.Done()
					cErr <- e
				}()

				err = <-cErr
				close(cErr)

				if err != nil {
					totalErrors++
					log.Printf("error to process %d", i)
					log.Println("--- checkValidData: ", data, " ---")
					err = fmt.Errorf("s.repository.InsertSanitizedData(): %w", err)
					log.Println(err)
					return
				}
			}
			wg.Wait()

			return
		}

		var sem = semaphore.NewWeighted(int64(10))

		values, e := s.repository.GetData(limit, offset)
		if err != nil {
			err = fmt.Errorf("s.repository.GetData(): %w", e)
			return
		}

		c1, c2 := make(chan int), make(chan error)

		if e := sem.Acquire(context.Background(), 1); e != nil {
			err = fmt.Errorf("sem.Acquire(): %w", e)
			return
		}

		go func() {

			totalErrors, e := sanitizeData(values, limit, offset)
			if err != nil {
				err = fmt.Errorf("sanitizeData(): %w", e)
			}
			c1 <- totalErrors
			c2 <- err

			sem.Release(1)

		}()

		totalErr := <-c1
		err = <-c2

		close(c1)
		close(c2)

		log.Printf("group %d has %d errors \n", page, totalErr)
	}

	ok = true
	err = nil

	return
}

//CountSanitizedData used to check reliability
func (s *InternalService) CountSanitizedData() (totalLines int, err error) {
	return s.repository.CountSanitizedData()
}
