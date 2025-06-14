package reporting

type Service interface {
	GenerateReport(data []byte) error
}
