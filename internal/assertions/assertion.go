package assertions

// BaseAssertion contains elements mandatory to all assertions.
type BaseAssertion struct {
	// screenAsserter is the ScreenAsserter that contains the rendered screenstate
	screenAsserter *ScreenAsserter

	// rowIndex is the index of the row in the screenstate that we want to assert on
	rowIndex int

	// Each assertion should update the row index in screenAsserter atmost once
	// Once it performs the update, it should set this flag to true
	ifUpdatedRowIndex bool
}

type Assertion interface {
	Run() error
	WrappedRun() bool
	GetRowUpdateCount() int
	UpdateRowIndex()
	GetType() string
}
