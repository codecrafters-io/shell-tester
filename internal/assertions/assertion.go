package assertions

type Assertion interface {
	Run(screenState [][]string, startRowIndex int) (processedRowCount int, err *AssertionError)
	Inspect() string
}
