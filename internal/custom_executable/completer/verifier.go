package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
)

type verifier struct {
	config               completer_configuration.CompleterConfiguration
	completerSecretValue string
	accumulatedLogs      []completer_configuration.VerifierLog
}

func (v *verifier) appendLog(logType completer_configuration.VerifierLogType, message string) {
	v.accumulatedLogs = append(
		v.accumulatedLogs,
		completer_configuration.VerifierLog{Type: logType, Message: message},
	)
}

func (v *verifier) logSuccess(message string) {
	v.appendLog(completer_configuration.VerifierLogTypeSuccess, message)
}

func (v *verifier) logError(message string) {
	v.appendLog(completer_configuration.VerifierLogTypeError, message)
}

func writeStartupFailureReport(completerSecretValue string, message string, exitCode int) {
	report := completer_configuration.VerifierErrorReport{
		Logs: []completer_configuration.VerifierLog{
			{Type: completer_configuration.VerifierLogTypeError, Message: message},
		},
	}
	encodedReport, err := json.Marshal(report)
	if err != nil {
		os.Exit(exitCode)
	}
	path := completer_configuration.GetVerifierErrorReportPath(completerSecretValue)
	_ = os.WriteFile(path, encodedReport, 0644)
	os.Exit(exitCode)
}

func (v *verifier) areExpectationsViolated() bool {
	return v.argvExpectationsViolated() || v.envExpectationsViolated()
}

func (v *verifier) argvExpectationsViolated() bool {
	if v.config.ExpectedArguments == nil {
		return false
	}
	if len(os.Args) != 4 {
		return true
	}
	exp := v.config.ExpectedArguments
	return os.Args[1] != exp.Argv1 ||
		os.Args[2] != exp.Argv2 ||
		os.Args[3] != exp.Argv3
}

func (v *verifier) envExpectationsViolated() bool {
	if v.config.ExpectedEnvVars == nil {
		return false
	}
	exp := v.config.ExpectedEnvVars
	return os.Getenv("COMP_LINE") != exp.CompLine ||
		os.Getenv("COMP_POINT") != exp.CompPoint
}

func (v *verifier) emitErrorReport() {
	if v.config.ExpectedArguments != nil {
		v.emitArgvVerificationErrorReport()
	}

	if v.config.ExpectedEnvVars != nil {
		v.emitEnvironmentVerificationErrorReport()
	}
}

func (v *verifier) emitArgvVerificationErrorReport() {
	exp := v.config.ExpectedArguments
	argc := len(os.Args) - 1

	if len(os.Args) < 4 {
		v.logError(fmt.Sprintf(
			"Expected argv[1] through argv[3] to be present, found only up to argv[%d]",
			argc,
		))
		return
	}
	if len(os.Args) > 4 {
		v.logError(fmt.Sprintf(
			"Expected only argv[1] through argv[3] to be present, found up to argv[%d]",
			argc,
		))
		return
	}
	v.logSuccess("Arguments are of length 4")

	argv1, argv2, argv3 := os.Args[1], os.Args[2], os.Args[3]

	if argv1 != exp.Argv1 {
		v.logError(fmt.Sprintf("Expected argv[1] to be %q, found %q", exp.Argv1, argv1))
	} else {
		v.logSuccess("Expected value of argv[1] found")
	}

	if argv2 != exp.Argv2 {
		v.logError(fmt.Sprintf("Expected argv[2] to be %q, found %q", exp.Argv2, argv2))
	} else {
		v.logSuccess("Expected value of argv[2] found")
	}

	if argv3 != exp.Argv3 {
		v.logError(fmt.Sprintf("Expected argv[3] to be %q, found %q", exp.Argv3, argv3))
	} else {
		v.logSuccess("Expected value of argv[3] found")
	}
}

func (v *verifier) emitEnvironmentVerificationErrorReport() {
	exp := v.config.ExpectedEnvVars

	receivedCompLine := os.Getenv("COMP_LINE")
	if receivedCompLine == exp.CompLine {
		v.logSuccess("Expected value of environment variable COMP_LINE found")
	} else {
		v.logError(fmt.Sprintf(
			"Expected environment variable COMP_LINE to be %q, found %q",
			exp.CompLine, receivedCompLine,
		))
	}

	receivedCompPoint := os.Getenv("COMP_POINT")
	if receivedCompPoint == exp.CompPoint {
		v.logSuccess("Expected value of environment variable COMP_POINT found")
	} else {
		v.logError(fmt.Sprintf(
			"Expected environment variable COMP_POINT to be %q, found %q",
			exp.CompPoint, receivedCompPoint,
		))
	}
}

func (v *verifier) saveErrorReportToDiskAndExit(exitCode int) {
	report := completer_configuration.VerifierErrorReport{Logs: v.accumulatedLogs}

	encodedReport, err := json.Marshal(report)
	if err != nil {
		os.Exit(exitCode)
	}

	reportPath := completer_configuration.GetVerifierErrorReportPath(v.completerSecretValue)

	if err := os.WriteFile(reportPath, encodedReport, 0644); err != nil {
		os.Exit(exitCode)
	}

	os.Exit(exitCode)
}
