package custom_executable

import "fmt"

func CreateSignaturePrinterExecutable(randomString, outputPath string) error {
	if len(randomString) == 0 {
		return fmt.Errorf("CodeCrafters Internal Error: randomString must be non-empty")
	}

	if len(randomString) > secretSlotByteLength(1) {
		return fmt.Errorf("CodeCrafters Internal Error: randomString exceeds %d bytes", secretSlotByteLength(1))
	}

	return prepareSecretPatchedExecutable("signature_printer", outputPath, []string{randomString})
}
