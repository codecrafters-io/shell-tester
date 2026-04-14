package utils

import (
	"fmt"
	"path/filepath"

	"github.com/codecrafters-io/tester-utils/executable"
)

func MustGetExecutablePathForCommand(command string) string {
	absPath, err := executable.ResolveAbsolutePath(command)
	if err != nil {
		panic(fmt.Sprintf(
			"Codecrafters Internal Error - Failed to resolve absolute path for command %s: %s",
			command,
			err,
		))
	}

	// The absolute path could be a symlink, so that must be resolved
	resolvedAbsolutePath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		panic(fmt.Sprintf(
			"Codecrafters Internal Error - Failed to resolve symlink for path %s: %s",
			absPath,
			err,
		))
	}

	return resolvedAbsolutePath
}
