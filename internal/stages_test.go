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
			StageSlugs:          []string{"cz2"},
			CodePath:            "./test_helpers/scenarios/wrong_output",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/wrong_output",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"no_command_fail": {
			StageSlugs:          []string{"cz2"},
			CodePath:            "./test_helpers/scenarios/no_output",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/no_output",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"escape_codes_pass": {
			StageSlugs:          []string{"cz2"},
			CodePath:            "./test_helpers/scenarios/escape_codes",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/escape_codes",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"exit_error_fail": {
			StageSlugs:          []string{"pn5"},
			CodePath:            "./test_helpers/scenarios/exit_error",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/exit_error",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"run_program_newline_in_args": {
			StageSlugs:          []string{"ip1"},
			CodePath:            "./test_helpers/scenarios/run_program_newline_in_args",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/run_program_newline_in_args",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"stage_1_only_prompt_pass": {
			StageSlugs:          []string{"oo8"},
			CodePath:            "./test_helpers/scenarios/stage_1_only_prompt",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/stage_1_only_prompt",
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
		"filename_completions_pass_bash": {
			StageSlugs:          []string{"zv2", "ue6", "lc6", "vs5", "no5", "jp8", "bf8"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/filename_completions/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"background_jobs_pass_bash": {
			StageSlugs:          []string{"af3", "at7", "jd6", "dk5"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/background_jobs/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"background_jobs_incorrect_output_format": {
			StageSlugs:          []string{"at7"},
			CodePath:            "./test_helpers/scenarios/background_jobs_incorrect_output_format",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/background_jobs_incorrect_output_format",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"background_jobs_incorrect_job_number": {
			StageSlugs:          []string{"at7"},
			CodePath:            "./test_helpers/scenarios/background_jobs_incorrect_job_number",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/background_jobs_incorrect_job_number",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"background_jobs_jobs_builtin_incorrect_marker": {
			StageSlugs:          []string{"dk5"},
			CodePath:            "./test_helpers/scenarios/background_jobs_jobs_builtin_incorrect_marker",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/background_jobs_jobs_builtin_incorrect_marker",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"background_jobs_jobs_builtin_incorrect_output_format": {
			StageSlugs:          []string{"dk5"},
			CodePath:            "./test_helpers/scenarios/background_jobs_jobs_builtin_incorrect_output_format",
			ExpectedExitCode:    1,
			StdoutFixturePath:   "./test_helpers/fixtures/background_jobs_jobs_builtin_incorrect_output_format",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"pipelines_pass_bash": {
			StageSlugs:          []string{"br6", "ny9", "xk3"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/pipelines/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"history_pass_bash": {
			StageSlugs:          []string{"bq4", "yf5", "ag6", "rh7", "vq0", "dm2"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/history/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"history_persistence_pass_bash": {
			StageSlugs:          []string{"za2", "in3", "sx3", "kz7", "zp4"},
			CodePath:            "./test_helpers/bash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/bash/history_persistence/pass",
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
		"filename_completions_pass_ash": {
			StageSlugs:          []string{"zv2", "ue6", "lc6", "vs5", "no5", "jp8", "bf8"},
			CodePath:            "./test_helpers/ash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/ash/filename_completions/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"pipelines_pass_ash": {
			StageSlugs:          []string{"br6", "ny9", "xk3"},
			CodePath:            "./test_helpers/ash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/ash/pipelines/pass",
			NormalizeOutputFunc: normalizeTesterOutput,
		},
		"history_pass_ash": {
			StageSlugs:          []string{"bq4", "yf5", "rh7", "vq0", "dm2"},
			CodePath:            "./test_helpers/ash",
			ExpectedExitCode:    0,
			StdoutFixturePath:   "./test_helpers/fixtures/ash/history/pass",
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
		"[your-program] <cwd>":            {regexp.MustCompile(`\[your-program\] .{0,4}/(workspace|workspaces|home|Users|app)/.*`)},
		"ls-la-output-line":               {regexp.MustCompile(`-rw-r--r-- .*`)},
		"PATH is now: <path>":             {regexp.MustCompile(`PATH is now: .*`)},
		"/tmp/":                           {regexp.MustCompile(`/var/folders/.*/.*/.*/`)},
		"[your-program] [JOB_NUM] PID":    {regexp.MustCompile(`\[your-program\].*\[\d+\] \d+`)},
		// For background_jobs_incorrect_output_format
		"[your-program] [JOB_NUM]PID":               {regexp.MustCompile(`\[your-program\].*\[\d+\]\d+`)},
		"[tester::#AT7] Received: \"[JOB_NUM]PID\"": {regexp.MustCompile(`\[tester::#AT7\].*Received:.*"\[\d+\]\d+"`)},
		"[tester::#AT7] REGEX (HINT)":               {regexp.MustCompile(`\[tester::#AT7\].*\(Hint:.*\)`)},
	}

	for replacement, regexes := range replacements {
		for _, regex := range regexes {
			testerOutput = regex.ReplaceAll(testerOutput, []byte(replacement))
		}
	}

	return testerOutput
}
