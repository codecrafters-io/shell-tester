package assertions

type AssertionError struct {
	ErrorRowIndex int // Will be -1 if the error doesn't affect a specific line range (e.g. bell)
	Message       string
}

func (e AssertionError) AffectsLine() bool {
	return e.ErrorRowIndex != -1
}

func (e AssertionError) Error() string {
	return `CodeCrafters Internal Error: AssertionError#Error() should not be called`
}

func (e AssertionError) ErrorMessage() string {
	return e.Message
}
