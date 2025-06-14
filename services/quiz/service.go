package quiz

type Service interface {
	ProcessQuiz(data []byte) error
}
