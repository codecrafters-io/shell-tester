package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
)

func main() {
	completerSecretValue := "<<RANDOM>>"

	configPath := filepath.Join("/tmp", completerSecretValue)
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		writeStartupFailureReport(
			completerSecretValue,
			fmt.Sprintf("Codecrafters Internal Error - Could not read completer configuration at %s: %v", configPath, err),
			1,
		)
	}

	var config completer_configuration.CompleterConfiguration
	if err := json.Unmarshal(configBytes, &config); err != nil {
		writeStartupFailureReport(
			completerSecretValue,
			fmt.Sprintf("Codecrafters Internal Error - Could not parse completer configuration JSON: %v", err),
			1,
		)
	}

	verifier := &verifier{
		config:               config,
		completerSecretValue: completerSecretValue,
	}

	if verifier.areExpectationsViolated() {
		verifier.emitErrorReport()
		verifier.saveErrorReportToDiskAndExit(1)
	}

	for _, outputLine := range config.OutputLines {
		fmt.Println(outputLine)
	}
}
