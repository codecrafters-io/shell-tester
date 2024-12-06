package assertions

import "github.com/codecrafters-io/tester-utils/logger"

type Assertion interface {
	Run(screenState [][]string, logger *logger.Logger) error
}
