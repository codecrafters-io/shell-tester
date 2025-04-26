package internal

import (
	"os"
	"regexp"
	"runtime"
	"testing"

	testerUtilsTesting "github.com/codecrafters-io/tester-utils/testing"
)

func TestStages(t *testing.T) {
	os.Setenv("CODECRAFTERS_RANDOM_SEED", "1234567890")

	testCases := map[string]testerUtilsTesting.TesterOutputTestCase{
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
		"base_pass_bash": {
			UntilStageSlug:      "ip1",
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/base/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"navigation_pass_bash": {
			UntilStageSlug:      "gp4",
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/navigation/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"quoting_pass_bash": {
			StageSlugs:          []string{"ni6", "tg6", "yt5", "le5", "gu3", "qj0"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/quoting/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"redirection_pass_bash": {
			StageSlugs:          []string{"jv1", "vz4", "el9", "un3"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/redirection/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"completions_pass_bash": {
			StageSlugs:          []string{"qp2", "gm9", "qm8", "gy5", "wh6", "wt6"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/completions/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"base_pass_ash": {
			UntilStageSlug:      "ip1",
			CodePath:            "./test_helpers/ash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/ash/base/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"navigation_pass_ash": {
			UntilStageSlug:      "gp4",
			CodePath:            "./test_helpers/ash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/ash/navigation/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"quoting_pass_ash": {
			StageSlugs:          []string{"ni6", "tg6", "yt5", "le5", "gu3", "qj0"},
			CodePath:            "./test_helpers/ash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/ash/quoting/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"redirection_pass_ash": {
			StageSlugs:          []string{"jv1", "vz4", "el9", "un3"},
			CodePath:            "./test_helpers/ash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/ash/redirection/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"completions_pass_ash": {
			// TODO debug why this stage fails in make test ?
			// "gy5"
			StageSlugs:          []string{"qp2", "gm9", "qm8", "wh6", "wt6"},
			CodePath:            "./test_helpers/ash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/ash/completions/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
	}

	if runtime.GOOS == "darwin" {
		// Getting almquist shell (ash) to work properly on macOS is a pain.
		// So, we skip those while running make test on macOS.
		testCases = filterOutAshTestCases(testCases)
	}
	testerUtilsTesting.TestTesterOutput(t, testerDefinition, testCases)
}

func filterOutAshTestCases(testCases map[string]testerUtilsTesting.TesterOutputTestCase) map[string]testerUtilsTesting.TesterOutputTestCase {
	filteredTestCases := make(map[string]testerUtilsTesting.TesterOutputTestCase)
	for slug, testCase := range testCases {
		if testCase.CodePath != "./test_helpers/ash" {
			filteredTestCases[slug] = testCase
		}
	}
	return filteredTestCases
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
