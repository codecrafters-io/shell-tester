package assertions

type AssertionError struct {
	StartRowIndex int
	ErrorRowIndex int
	Message       string
}

func (e AssertionError) Error() string {
	return `CodeCrafters Internal Error: AssertionError#Error() should not be called`
}

func (e AssertionError) ErrorMessage() string {
	return e.Message
}
