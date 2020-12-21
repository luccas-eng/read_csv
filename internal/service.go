package internal

//Service interface implement services running between clients
type Service interface {
}

//InternalService struct implements repo
type InternalService struct {
	repository Repository
}

//NewService implements the repo service
func NewService(r Repository) Service {
	return &InternalService{r}
}
