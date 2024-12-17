package internal

import (
	"os"
	"regexp"
	"testing"

	testerUtilsTesting "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	os.Setenv("CODECRAFTERS_RANDOM_SEED", "1234567890")

	testCases := map[string]testerUtilsTesting.TesterOutputTestCase{
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
		"missing_command_fail": {
			UntilStageSlug:      "cz2",
			CodePath:            "./test_helpers/scenarios/wrong_output",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/wrong_output",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"no_command_fail": {
			UntilStageSlug:      "cz2",
			CodePath:            "./test_helpers/scenarios/no_output",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/no_output",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"escape_codes_pass": {
			UntilStageSlug:      "cz2",
			CodePath:            "./test_helpers/scenarios/escape_codes",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/escape_codes",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"quoting_pass": {
			StageSlugs:          []string{"ni6", "tg6", "yt5", "le5", "gu3", "qj0"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/quoting/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	testerUtilsTesting.TestTesterOutput(t, testerDefinition, testCases)
}

func normalizeTesterOutput(testerOutput []byte) []byte {
	replacements := map[string][]*regexp.Regexp{
		"/bin/$1":                         {regexp.MustCompile(`\/usr/bin/(\w+)`)},
		"[your-program] my_exe is <path>": {regexp.MustCompile(`\[your-program\] .{4}my_exe is .*`)},
		"[your-program] <cwd>":            {regexp.MustCompile(`\[your-program\] .{4}/(workspaces|home|Users)/.*`)},
	}

	for replacement, regexes := range replacements {
		for _, regex := range regexes {
			testerOutput = regex.ReplaceAll(testerOutput, []byte(replacement))
		}
	}

	return testerOutput
}
