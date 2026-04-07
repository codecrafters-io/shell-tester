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
