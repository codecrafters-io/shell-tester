package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
)

type verifier struct {
	config                      completer_configuration.CompleterConfiguration
	accumulatedVerificationLogs []string
}

func (v *verifier) RegisterLog(line string) {
	v.accumulatedVerificationLogs = append(v.accumulatedVerificationLogs, line)
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

func (v *verifier) verificationErrorReport() string {
	v.produceVerificationFailureLogs()
	return strings.Join(v.accumulatedVerificationLogs, "\n")
}

func (v *verifier) produceVerificationFailureLogs() {
	if v.config.ExpectedArguments != nil {
		v.emitArgvVerificationSection(v.config)
	}

	if v.config.ExpectedEnvVars != nil {
		v.emitEnvironmentVerificationSection(v.config)
	}
}

func (v *verifier) emitArgvVerificationSection(config completer_configuration.CompleterConfiguration) {
	exp := config.ExpectedArguments
	argc := len(os.Args) - 1

	if len(os.Args) < 4 {
		v.RegisterLog(fmt.Sprintf(
			"☓ Expected argv[1] through argv[3] to be present, found only up to argv[%d]",
			argc,
		))
		return
	} else if len(os.Args) > 4 {
		v.RegisterLog(fmt.Sprintf(
			"☓ Expected only argv[1] through argv[3] to be present, found up to argv[%d]",
			argc,
		))
		return
	} else {
		v.RegisterLog("✓ Arguments are of length 4")
	}

	argv1, argv2, argv3 := os.Args[1], os.Args[2], os.Args[3]

	if argv1 != exp.Argv1 {
		v.RegisterLog(fmt.Sprintf("☓ Expected argv[1] to be %q, found %q", exp.Argv1, argv1))
	} else {
		v.RegisterLog("✓ Expected value of argv[1] found")
	}

	if argv2 != exp.Argv2 {
		v.RegisterLog(fmt.Sprintf("☓ Expected argv[2] to be %q, found %q", exp.Argv2, argv2))
	} else {
		v.RegisterLog("✓ Expected value of argv[2] found")
	}

	if argv3 != exp.Argv3 {
		v.RegisterLog(fmt.Sprintf("☓ Expected argv[3] to be %q, found %q", exp.Argv3, argv3))
	} else {
		v.RegisterLog("✓ Expected value of argv[3] found")
	}
}

func (v *verifier) emitEnvironmentVerificationSection(config completer_configuration.CompleterConfiguration) {
	exp := config.ExpectedEnvVars

	receivedCompLine := os.Getenv("COMP_LINE")
	if receivedCompLine == exp.CompLine {
		v.RegisterLog("✓ Expected value of environment variable COMP_LINE found")
	} else {
		v.RegisterLog(fmt.Sprintf(
			"☓ Expected environment variable COMP_LINE to be %q found %q",
			exp.CompLine, receivedCompLine,
		))
	}

	receivedCompPoint := os.Getenv("COMP_POINT")
	if receivedCompPoint == exp.CompPoint {
		v.RegisterLog("✓ Expected value of environment variable COMP_POINT found")
	} else {
		v.RegisterLog(fmt.Sprintf(
			"☓ Expected environment variable COMP_POINT to be %q found %q",
			exp.CompPoint, receivedCompPoint,
		))
	}
}

func main() {
	secretCode := "<<RANDOM>>"

	configPath := filepath.Join("/tmp", secretCode)
	data, err := os.ReadFile(configPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	var config completer_configuration.CompleterConfiguration

	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Verify args and environment
	verifier := &verifier{
		config: config,
	}

	if verifier.areExpectationsViolated() {
		fmt.Fprintln(
			os.Stderr,
			completerScriptError(verifier.verificationErrorReport()),
		)
		os.Exit(1)
	}

	streamOut := os.Stdout
	if config.UseStderrStream {
		streamOut = os.Stderr
	}

	for _, outputLine := range config.OutputLines {
		fmt.Fprintln(streamOut, outputLine)
	}

	// Sleep for 120ms
	// Completer which is bound to generate logs to stderr should never exit
	// to test that the shell has actually streamed the stderr instead of collecting it
	// after it has exitted
	if config.UseStderrStream {
		time.Sleep(120 * time.Second)
	}
}

func completerScriptError(message string) string {
	return "\nError from the completer script:\n" + message
}
