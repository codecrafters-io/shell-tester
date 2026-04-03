package custom_executable

import (
	"fmt"
	"strings"
)

// CreateLCPEnvCompleter builds a completer that checks argv[1]/argv[3], validates COMP_LINE/COMP_POINT consistency
// (cursor at end-of-line), then prints newline-separated candidates filtered by argv[2] prefix (checkout, cherry-pick).
func CreateLCPEnvCompleter(outputPath, envLineVar, envPointVar, command, prevWord string) error {
	vals := []string{envLineVar, envPointVar, command, prevWord}
	secrets := make([]string, len(vals))
	for i, s := range vals {
		slot := i + 1
		max := SecretSlotByteLen(slot)
		if len(s) > max {
			return fmt.Errorf("CodeCrafters Internal Error: parameter %d exceeds %d bytes", i+1, max)
		}
		secrets[i] = s + strings.Repeat(" ", max-len(s))
	}
	return prepareSecretPatchedExecutable(secretPatchedLCPEnvCompleter, outputPath, secrets)
}
