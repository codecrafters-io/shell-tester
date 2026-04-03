package custom_executable

import (
	"fmt"
	"strings"
)

// CreateMultiCandidateEnvCompleter builds a completer that validates argv and COMP_LINE / COMP_POINT (via patched env var names),
// then prints newline-separated completion candidates from a fixed list filtered by prefix match on argv[2].
func CreateMultiCandidateEnvCompleter(outputPath, envLineVar, envPointVar, arg1, arg2, arg3, wantCompLine, wantCompPoint string) error {
	vals := []string{envLineVar, envPointVar, arg1, arg2, arg3, wantCompLine, wantCompPoint}
	secrets := make([]string, len(vals))
	for i, s := range vals {
		slot := i + 1
		max := SecretSlotByteLen(slot)
		if len(s) > max {
			return fmt.Errorf("CodeCrafters Internal Error: parameter %d exceeds %d bytes", i+1, max)
		}
		secrets[i] = s + strings.Repeat(" ", max-len(s))
	}
	return prepareSecretPatchedExecutable(secretPatchedMultiCandidateEnvCompleter, outputPath, secrets)
}
