package internal

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

//Service interface implement services running between clients
type Service interface {
	ReadData() error
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
func (s *InternalService) ReadData() error {
	data, err := ioutil.ReadFile("./external/base_teste.txt")
	if err != nil {
		return fmt.Errorf("ioutil.ReadFile(): %w", err)
	}

	r := regexp.MustCompile("[^\\s]+")
	values := r.FindAllString(string(data), -1)
	for _, i := range values {
		fmt.Println(i)
	}
	return nil
}
