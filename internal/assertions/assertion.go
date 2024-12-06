package assertions

type Assertion interface {
	Run(value string) error
}
