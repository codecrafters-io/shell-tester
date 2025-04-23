package custom_executable

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

// func CreateSignaturePrinterExecutable(randomString, outputPath string) error {
// 	if len(randomString) != 10 {
// 		return fmt.Errorf("CodeCrafters Internal Error: randomString length must be 10")
// 	}

// 	// Copy the base executable from archive location to user's executable path
// 	err := copyExecutable("signature_printer", outputPath)
// 	if err != nil {
// 		return fmt.Errorf("CodeCrafters Internal Error: copying executable failed: %w", err)
// 	}

// 	// Replace the placeholder with the random string
// 	// We can run the executable now, it will work as expected
// 	err = addSecretCodeToExecutable(outputPath, randomString)
// 	if err != nil {
// 		return fmt.Errorf("CodeCrafters Internal Error: adding secret code to executable failed: %w", err)
// 	}

// 	return nil
// }

const secretCodeVariablePath = "main.secretCode"

func CreateSignaturePrinterExecutable(randomString, outputPath string) error {
	if len(randomString) != 10 {
		return fmt.Errorf("CodeCrafters Internal Error: randomString length must be 10")
	}

	oses := []string{"darwin", "linux"}
	arches := []string{"arm64", "amd64"}

	ldflags := fmt.Sprintf("-X '%s=%s'", secretCodeVariablePath, randomString)

	name := "signature_printer"
	sourcePath := path.Join(os.Getenv("TESTER_DIR"), "internal", "custom_executable", "signature_printer", "main.go")

	for _, goos := range oses {
		for _, goarch := range arches {
			outputPath := path.Join(os.Getenv("TESTER_DIR"), "built_executables", fmt.Sprintf("%s_%s_%s", name, goos, goarch))
			fmt.Printf("Building %s for %s/%s -> %s\n", name, goos, goarch, outputPath)

			// Prepare the command
			cmd := exec.Command("go", "build", "-o", outputPath, "-ldflags", ldflags, sourcePath)

			// Set environment variables for cross-compilation
			cmd.Env = append(os.Environ(),
				fmt.Sprintf("GOOS=%s", goos),
				fmt.Sprintf("GOARCH=%s", goarch),
			)

			// Run the command and capture output
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Error building %s for %s/%s: %v\nOutput:\n%s", name, goos, goarch, err, string(output))
				// Decide if you want to stop on error or continue
				os.Exit(1) // uncomment to stop on first error
			} else if len(output) > 0 {
				// Print successful build output if any (usually none unless verbose)
				fmt.Printf("Output for %s %s/%s:\n%s\n", name, goos, goarch, string(output))
			}
		}
	}

	fmt.Println("Build process completed.")

	return nil
}
