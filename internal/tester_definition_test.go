package internal

import (
	"testing"

	testerUtilsTesting "github.com/codecrafters-io/tester-utils/testing"
)

func TestStagesMatchYAML(t *testing.T) {
	testerUtilsTesting.ValidateTesterDefinitionAgainstYAML(t, testerDefinition, "test_helpers/course_definition.yml")
}
