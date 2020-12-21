package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/read_csv/internal/model"
)

//Service interface implement services running between clients
type Service interface {
	ProcessData(c context.Context) (values []*model.MapData, total int, err error)
}

//InternalService struct implements repo
type InternalService struct {
	repository Repository
}

//NewService implements the repo service
func NewService(r Repository) Service {
	return &InternalService{r}
}

//ReadData ...
func (s *InternalService) ProcessData(c context.Context) (values []*model.MapData, total int, err error) {

	file, err := os.Open("../external/base_teste.txt")
	if err != nil {
		return nil, 0, fmt.Errorf("os.Open(): %w", err)
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
					return nil, 0, fmt.Errorf("%w", err)
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
			return nil, 0, fmt.Errorf("s.repository.InsertValues(): %w", err)
		}

		data := &model.MapData{Key: counter, Value: v}

		values = append(values, data)

		counter++

	}

	if err != nil && err != io.EOF {
		return nil, 0, fmt.Errorf("process failed with error: %w", err)
	}

	total = counter

	return
}
