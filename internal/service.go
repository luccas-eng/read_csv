package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/miguelpragier/handy"
)

//Service interface implement services running between clients
type Service interface {
	ProcessData(filePath string) (total int, err error)
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

const mb = 1024 * 1024
const gb = 1024 * mb

//ProcessData ...
func (s *InternalService) ProcessData(filePath string) (total int, err error) {

	dat, _ := ioutil.ReadFile(filePath)
	sdat := bytes.Split(dat, []byte{'\n'})

	for i := range sdat {

		if ok, err := regexp.MatchString("CPF", string(sdat[i])); ok && err == nil {
			continue
		}

		fields := bytes.Fields(sdat[i])

		data, e := sanitizeData(fields)
		if e != nil {
			err = fmt.Errorf("sanitizeData(): %w", e)
			return
		}

		err = s.repository.InsertSanitizedData(data)
		if err != nil {
			return 0, fmt.Errorf("s.repository.InsertValues(): %w", err)
		}
	}

	return
}

func sanitizeData(values [][]byte) (data []interface{}, err error) {

	var invalidCPF, invalidCNPJ bool

	if len(values) < 8 {
		err = fmt.Errorf("len(values) < 8")
		return
	}

	for f := range values {
		values[f] = bytes.ToLower(bytes.TrimSpace(values[f]))
	}

	if !bytes.Equal(values[0], []byte("null")) {
		if !handy.CheckCPF(string(values[0])) {
			invalidCPF = true
			data = append(data, handy.OnlyLettersAndNumbers(string(values[0])))
		} else {
			data = append(data, handy.OnlyDigits(string(values[0])))
		}
	} else {
		data = append(data, "")
	}

	if !bytes.Equal(values[1], []byte("null")) {
		data = append(data, handy.OnlyLettersAndNumbers(string(values[1])))
	} else {
		data = append(data, "")
	}

	if !bytes.Equal(values[2], []byte("null")) {
		data = append(data, handy.OnlyLettersAndNumbers(string(values[2])))
	} else {
		data = append(data, "")
	}

	if !bytes.Equal(values[3], []byte("null")) {
		data = append(data, handy.OnlyLettersAndNumbers(string(values[3])))
	} else {
		data = append(data, "")
	}

	if !bytes.Equal(values[4], []byte("null")) {
		avgTicket, e := strconv.ParseFloat(strings.ReplaceAll(string(values[4]), ",", "."), 64)
		if e != nil {
			err = fmt.Errorf("strconv.ParseFloat(): %w", e)
			return
		}
		data = append(data, avgTicket)
	} else {
		data = append(data, 0.0)
	}

	if !bytes.Equal(values[5], []byte("null")) {
		lastTicket, e := strconv.ParseFloat(strings.ReplaceAll(string(values[5]), ",", "."), 64)
		if e != nil {
			err = fmt.Errorf("strconv.ParseFloat(): %w", e)
			return
		}
		data = append(data, lastTicket)
	} else {
		data = append(data, 0.0)
	}

	if !bytes.Equal(values[6], []byte("null")) {
		if !handy.CheckCNPJ(string(values[6])) {
			invalidCNPJ = true
			data = append(data, handy.OnlyLettersAndNumbers(string(values[6])))
		} else {
			data = append(data, handy.OnlyDigits(string(values[6])))
		}
	} else {
		data = append(data, "")
	}

	if !bytes.Equal(values[7], []byte("null")) {
		if !handy.CheckCNPJ(string(values[7])) {
			invalidCNPJ = true
			data = append(data, handy.OnlyLettersAndNumbers(string(values[7])))
		} else {
			data = append(data, handy.OnlyDigits(string(values[7])))
		}
	} else {
		data = append(data, "")
	}

	if invalidCPF {
		data = append(data, true)
	}

	if invalidCNPJ {
		data = append(data, true)
	}

	if err != nil {
		log.Println("invalidData: ", data)
		log.Println(err)
		return
	}

	return
}

//CountSanitizedData used to check reliability
func (s *InternalService) CountSanitizedData() (totalLines int, err error) {
	return s.repository.CountSanitizedData()
}
