package assertions

type Assertion interface {
	Run() error
	WrappedRun() bool
	GetRowUpdateCount() int
	UpdateRowIndex()
}
