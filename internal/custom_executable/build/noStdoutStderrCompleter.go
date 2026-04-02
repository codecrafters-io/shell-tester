package custom_executable

import (
	"fmt"
	"strings"
)

// CreateNoCompleterWithStderr builds a completer that prints nothing to stdout and writes exactly three
// stderr lines (bytes preserved, including no extra content). stderr must contain three newline-separated
// lines after trimming a single trailing newline; each line is at most SecretSlotByteLen(n) bytes.
func CreateNoCompleterWithStderr(outputPath, stderr string) error {
	stderr = strings.TrimSuffix(stderr, "\n")
	lines := strings.Split(stderr, "\n")
	if len(lines) != 3 {
		return fmt.Errorf("CodeCrafters Internal Error: stderr must contain exactly 3 lines separated by newlines, got %d line(s)", len(lines))
	}
	secrets := make([]string, 3)
	for i, line := range lines {
		slot := i + 1
		max := SecretSlotByteLen(slot)
		if len(line) > max {
			return fmt.Errorf("CodeCrafters Internal Error: stderr line %d exceeds %d bytes", i+1, max)
		}
		secrets[i] = line + strings.Repeat(" ", max-len(line))
	}
	return prepareSecretPatchedExecutable(secretPatchedNoStdoutStderrCompleter, outputPath, secrets)
}
