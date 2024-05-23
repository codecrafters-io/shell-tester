package internal

import (
	"os"
	"regexp"
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	os.Setenv("CODECRAFTERS_RANDOM_SEED", "1234567890")

	testCases := map[string]tester_utils_testing.TesterOutputTestCase{
		"base_stages_pass": {
			UntilStageSlug:      "ip1",
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/base/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"navigation_pass": {
			UntilStageSlug:      "gp4",
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/navigation/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	tester_utils_testing.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	replacements := map[string][]*regexp.Regexp{}

	for replacement, regexes := range replacements {
		for _, regex := range regexes {
			testerOutput = regex.ReplaceAll(testerOutput, []byte(replacement))
		}
	}

	return testerOutput
}
