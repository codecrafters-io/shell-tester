package custom_executable

// CreateNoCompleter builds a completer that prints nothing to stdout and exits successfully.
func CreateNoCompleter(outputPath string) error {
	if err := createExecutableForOSAndArch("no_stdout_completer", outputPath); err != nil {
		return err
	}
	return reSignExecutableDarwinARM64(outputPath)
}
