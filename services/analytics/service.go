package analytics

type Service interface {
	ProcessAnalytics(data []byte) error
}
