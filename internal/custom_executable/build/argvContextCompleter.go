package custom_executable

// CreateArgvContextCompleter builds the PA6 completer that validates argv[1]–argv[3] and prints a fixed completion on success.
func CreateArgvContextCompleter(outputPath string) error {
	if err := createExecutableForOSAndArch("argv_context_completer", outputPath); err != nil {
		return err
	}
	return reSignExecutableDarwinARM64(outputPath)
}
