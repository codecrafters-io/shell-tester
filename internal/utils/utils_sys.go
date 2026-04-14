package utils

import (
	"fmt"
	"path/filepath"

	"github.com/codecrafters-io/tester-utils/executable"
)

// MustGetExecutablePathAndResolvedSymlinkForCommand returns the absolute path for the command's executable
// followed by the resolved path of the symlink if the absolute path is a symlink
// For example, For 'sleep' if the absolute path is "/bin/sleep" which is a symbolic link
// to "/bin/busybox", it returns ("/bin/sleep", "/bin/busybox")
// For paths which are not symlinks, same value is returned twice
func MustGetExecutablePathAndResolvedSymlinkForCommand(command string) (string, string) {
	absolutePath, err := executable.ResolveAbsolutePath(command)
	if err != nil {
		panic(fmt.Sprintf(
			"Codecrafters Internal Error - Failed to resolve absolute path for command %s: %s",
			command,
			err,
		))
	}

	// The absolute path could be a symlink, so that must be resolved
	resolvedAbsolutePath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		panic(fmt.Sprintf(
			"Codecrafters Internal Error - Failed to resolve symlink for path %s: %s",
			absolutePath,
			err,
		))
	}

	return absolutePath, resolvedAbsolutePath
}
