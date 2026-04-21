package completer_configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/tester-utils/logger"
)

// VerifierLogType classifies a line in the completer verification report for tester logging.
type VerifierLogType string

const (
	VerifierLogTypeSuccess VerifierLogType = "SUCCESS"
	VerifierLogTypeError   VerifierLogType = "ERROR"
)

// VerifierLog is one recorded line from the completer executable's self-check.
type VerifierLog struct {
	Type    VerifierLogType
	Message string
}

// VerifierErrorReport is serialized to disk by the completer when argv/env expectations fail
// or when the completer cannot load its JSON configuration.
type VerifierErrorReport struct {
	Logs []VerifierLog
}

type CompleterConfigurationExpectedArguments struct {
	Argv1 string
	Argv2 string
	Argv3 string
}

type CompleterConfigurationExpectedEnvVars struct {
	CompLine  string
	CompPoint string
}

type CompleterConfiguration struct {
	OutputLines       []string
	ExpectedArguments *CompleterConfigurationExpectedArguments
	ExpectedEnvVars   *CompleterConfigurationExpectedEnvVars
}

// GetVerifierErrorReportPath is where the completer binary writes error reprot in json format
// for the tester to read after a failing completion assertion.
func GetVerifierErrorReportPath(completerSecretValue string) string {
	return filepath.Join("/tmp", fmt.Sprintf("verifier_report.%s.json", completerSecretValue))
}

// LogCompleterErrors reads the report file produced by the completer executable and logs each
// entry with the appropriate logger level (success / info / error).
func LogCompleterErrors(logger *logger.Logger, completerSecretValue string) {
	verifierErrorReportPath := GetVerifierErrorReportPath(completerSecretValue)
	errorReportBytes, err := os.ReadFile(verifierErrorReportPath)

	// It's fine if there are no errors found
	// One case is when the user might have called the completer script properly (producing no errors in the disk)
	// and have not used those candidates properly
	// In those cases, report won't be present in the disk
	if err != nil {
		return
	}

	var report VerifierErrorReport

	if err := json.Unmarshal(errorReportBytes, &report); err != nil {
		panic(
			fmt.Sprintf(
				"Codecrafters Internal Error - Failed to parse completer verifier report from %s: %v",
				verifierErrorReportPath,
				err,
			),
		)
	}

	logger.Infof("Errors from the completer script:")
	for _, entry := range report.Logs {
		switch entry.Type {
		case VerifierLogTypeSuccess:
			logger.Successf("✓ %s", entry.Message)
		case VerifierLogTypeError:
			logger.Errorf("%s", entry.Message)
		default:
			panic(fmt.Sprintf("Codecrafters Internal Error - Unknown type for VerifierLogType: %s", entry.Type))
		}
	}
}
