package ingestion

type Service interface {
	ProcessEvent(data []byte) error
}

type service struct {
}

func (s *service) ProcessEvent(data []byte) error {
	return nil
}

func New() Service {
	return &service{}
}
