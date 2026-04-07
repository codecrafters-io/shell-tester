package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/codecrafters-io/shell-tester/internal/custom_executable/completer/completer_configuration"
)

func main() {
	secretCode := "<<RANDOM>>"

	configPath := filepath.Join("/tmp", secretCode)
	data, err := os.ReadFile(configPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	var cfg completer_configuration.CompleterConfiguration

	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// This will almost always never occur because the configuration is always verified
	// while being written
	// This check is here just to stay on the safe side and to ensure that the
	// configuration file is not tampered with in between writing and reading
	if err := cfg.Verify(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := verifyArgs(&cfg); err != nil {
		fmt.Fprintln(os.Stderr, completerScriptError(err))
		os.Exit(1)
	}

	if err := verifyEnvVars(&cfg); err != nil {
		fmt.Fprintln(os.Stderr, completerScriptError(err))
		os.Exit(1)
	}

	// If there are specified stderrLines, print them and sleep for 120 seconds
	// and exit
	if len(cfg.StderrLines) > 0 {
		for _, line := range cfg.StderrLines {
			fmt.Fprintln(os.Stderr, line)
		}
		time.Sleep(120 * time.Millisecond)
		os.Exit(1)
	}

	// Print the completion candidates
	for _, line := range cfg.CompletionCandidates {
		fmt.Println(line)
	}
}

func completerScriptError(err error) string {
	return "\nError from the completer script:\n" + err.Error()
}

func verifyArgs(cfg *completer_configuration.CompleterConfiguration) error {
	if cfg.ExpectedArguments == nil {
		return nil
	}

	exp := cfg.ExpectedArguments
	n := len(os.Args) - 1
	if n < 3 {
		return fmt.Errorf("Expected args #1 thru #3 to be present, only found #%d", n)
	}

	checks := []struct {
		idx  int
		want string
		got  string
	}{
		{1, exp.Argv1, os.Args[1]},
		{2, exp.Argv2, os.Args[2]},
		{3, exp.Argv3, os.Args[3]},
	}
	for _, c := range checks {
		if c.got != c.want {
			return fmt.Errorf("Expected argv[%d] to be %q found %q", c.idx, c.want, c.got)
		}
	}

	return nil
}

func verifyEnvVars(cfg *completer_configuration.CompleterConfiguration) error {
	if cfg.ExpectedEnvVars == nil {
		return nil
	}

	e := cfg.ExpectedEnvVars
	checks := []struct {
		name string
		want string
		got  string
	}{
		{"COMP_LINE", e.CompLine, os.Getenv("COMP_LINE")},
		{"COMP_POINT", e.CompPoint, os.Getenv("COMP_POINT")},
	}
	for _, c := range checks {
		if c.got != c.want {
			return fmt.Errorf("Expected environment variable %s to be %q found %q", c.name, c.want, c.got)
		}
	}

	return nil
}
