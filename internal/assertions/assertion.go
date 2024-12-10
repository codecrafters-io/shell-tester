package assertions

type AssertionError struct {
	StartRowIndex int
	ErrorRowIndex int
	Message       string
}

type Assertion interface {
	Run(screenState [][]string, startRowIndex int) (processedRowCount int, err error)
}
