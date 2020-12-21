package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
)

//Service interface implement services running between clients
type Service interface {
	ReadData() ([]string, error)
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
func (s *InternalService) ReadData() (values []string, err error) {

	file, err := os.Open("../external/base_teste.txt")
	if err != nil {
		return nil, fmt.Errorf("os.Open(): %w", err)
	}

	defer file.Close()

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
					return nil, fmt.Errorf("%w", err)
				}
				break
			}
		}

		if err == io.EOF {
			break
		}

		//Discard first line
		if ok, err := regexp.MatchString("CPF", buffer.String()); ok && err == nil {
			continue
		}

		r := regexp.MustCompile("[^\\s]+")

		values = r.FindAllString(buffer.String(), -1)
		err = s.repository.InsertValues(values)
		if err != nil {
			return nil, fmt.Errorf("s.repository.InsertValues(): %w", err)
		}
	}

	if err != io.EOF {
		return nil, fmt.Errorf(" > Failed with error: %w", err)
	}

	return values, nil
}
