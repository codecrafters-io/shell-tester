package custom_executable

import (
	"fmt"
	"strings"
)

// CreateSingleCompleterExecutable copies the single_completer binary to outputPath and patches
// <<RANDOM_1>> so the program prints exactly subcommand (one line on stdout).
// subcommand must be non-empty and at most SecretSlotByteLen(1) bytes (space-padded), matching patch width.
func CreateSingleCompleterExecutable(outputPath, subcommand string) error {
	maxLen := SecretSlotByteLen(1)
	if len(subcommand) == 0 || len(subcommand) > maxLen {
		return fmt.Errorf("CodeCrafters Internal Error: subcommand length must be 1..%d", maxLen)
	}

	padded := subcommand + strings.Repeat(" ", maxLen-len(subcommand))
	return prepareSecretPatchedExecutable(secretPatchedSingleCompleter, outputPath, []string{padded})
}
