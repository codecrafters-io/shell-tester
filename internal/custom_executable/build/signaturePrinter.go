package custom_executable

import (
	"fmt"
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

	sourcePath := path.Join(os.Getenv("TESTER_DIR"), "internal", "custom_executable", "signature_printer", "main.go")

	cmd := exec.Command("go", "build", "-o", outputPath, "-ldflags", ldflags, sourcePath)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic("CodeCrafters Internal Error: failed to build signature printer executable for " + goos + "/" + goarch + "\n" + string(output))
	}

	return nil
}
