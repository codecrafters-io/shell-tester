package internal

import (
	"path"

	custom_executable "github.com/codecrafters-io/shell-tester/internal/custom_executable/build"
	"github.com/codecrafters-io/shell-tester/internal/shell_executable"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// TODO: Move randomDir methods to tester utils
// TODO: Move SetupCustomCommands to shell_executable
// TODO: Logging is currently in stages, but it should be here
// TODO: How to log when multiple executables are set up?

type CommandDetails struct {
	// CommandType is the type of the command, e.g. "ls", "cat"
	CommandType string
	// CommandName is the name of the generated executable, e.g. "custom_exe_1234"
	CommandName string
	// CommandMetadata is any other metadata required for generating the command
	// SignaturePrinter: random code
	// Cat, Ls, Head, Tail, Wc, Yes: nothing
	CommandMetadata string
}

func SetUpCustomCommands(stageHarness *test_case_harness.TestCaseHarness, shell *shell_executable.ShellExecutable, commands []CommandDetails, useShorterDirectory bool) (string, error) {
	stageHarness.Logger.UpdateSecondaryPrefix("setup")
	createExecutableDirFunc := GetRandomDirectory
	if useShorterDirectory {
		createExecutableDirFunc = GetShortRandomDirectory
	}

	executableDir, err := createExecutableDirFunc(stageHarness)
	if err != nil {
		return "", err
	}
	// Add the random directory to PATH
	// (where the custom executable is copied to)
	stageHarness.Logger.Infof("export PATH=%s:$PATH", executableDir)
	shell.AddToPath(executableDir)

	for _, commandDetail := range commands {
		switch commandDetail.CommandType {
		case "ls":
			customLsPath := path.Join(executableDir, commandDetail.CommandName)
			err = custom_executable.CreateLsExecutable(customLsPath)
			if err != nil {
				return "", err
			}
		case "cat":
			customCatPath := path.Join(executableDir, commandDetail.CommandName)
			err = custom_executable.CreateCatExecutable(customCatPath)
			if err != nil {
				return "", err
			}
		case "signature_printer":
			customSignaturePrinterPath := path.Join(executableDir, commandDetail.CommandName)
			err = custom_executable.CreateSignaturePrinterExecutable(commandDetail.CommandMetadata, customSignaturePrinterPath)
			if err != nil {
				return "", err
			}
		}
	}
	stageHarness.Logger.ResetSecondaryPrefix()

	return executableDir, nil
}
