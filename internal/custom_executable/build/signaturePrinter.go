package custom_executable

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
)

const secretCodeVariablePath = "main.secretCode"

func CreateSignaturePrinterExecutable(randomString, outputPath string) error {
	if len(randomString) != 10 {
		return fmt.Errorf("CodeCrafters Internal Error: randomString length must be 10")
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH

	ldflags := fmt.Sprintf("-X '%s=%s'", secretCodeVariablePath, randomString)

	name := "signature_printer"
	sourcePath := path.Join(os.Getenv("TESTER_DIR"), "internal", "custom_executable", "signature_printer", "main.go")

	fmt.Printf("Building %s for %s/%s -> %s\n", name, goos, goarch, outputPath)
	cmd := exec.Command("go", "build", "-o", outputPath, "-ldflags", ldflags, sourcePath)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error building %s for %s/%s: %v\nOutput:\n%s", name, goos, goarch, err, string(output))
		panic("CodeCrafters Internal Error: failed to build signature printer executable for " + goos + "/" + goarch)
	} else if len(output) > 0 {
		fmt.Printf("Output for %s %s/%s:\n%s\n", name, goos, goarch, string(output))
	}

	return nil
}
