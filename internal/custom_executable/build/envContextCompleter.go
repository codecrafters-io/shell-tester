package custom_executable

import (
	"fmt"
	"strings"
)

// CreateEnvContextCompleter builds a completer that checks argv[1]–argv[3] and env vars named by envLineVar / envPointVar
// (e.g. COMP_LINE, COMP_POINT) against the expected values. On success it prints completionWord (one line on stdout).
// Each string must fit in the corresponding <<RANDOM_n>> slot (space-padded in the binary).
func CreateEnvContextCompleter(outputPath, envLineVar, envPointVar, arg1, arg2, arg3, wantCompLine, wantCompPoint, completionWord string) error {
	vals := []string{envLineVar, envPointVar, arg1, arg2, arg3, wantCompLine, wantCompPoint, completionWord}
	secrets := make([]string, len(vals))
	for i, s := range vals {
		slot := i + 1
		max := SecretSlotByteLen(slot)
		if len(s) > max {
			return fmt.Errorf("CodeCrafters Internal Error: parameter %d exceeds %d bytes", i+1, max)
		}
		secrets[i] = s + strings.Repeat(" ", max-len(s))
	}
	return prepareSecretPatchedExecutable(secretPatchedEnvContextCompleter, outputPath, secrets)
}
