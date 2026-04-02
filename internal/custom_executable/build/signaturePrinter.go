package custom_executable

import "fmt"

func CreateSignaturePrinterExecutable(randomString, outputPath string) error {
	// Embedded token is <<RANDOM_1>>; replacement must match its byte length exactly.
	if len(randomString) != SecretSlotByteLen(1) {
		return fmt.Errorf("CodeCrafters Internal Error: randomString length must be %d", SecretSlotByteLen(1))
	}

	return prepareSecretPatchedExecutable(secretPatchedSignaturePrinter, outputPath, []string{randomString})
}
