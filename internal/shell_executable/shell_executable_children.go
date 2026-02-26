package shell_executable

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

// GetChildPidFromCmdLine returns the PID of the first child process of the shell
// that has the same launch command as the provided string
// If it cannot find such a children process, it'll return -1
func (e *ShellExecutable) GetChildPidFromCmdLine(cmdLine string) int {
	shellPid := e.GetShellPid()
	processes, err := process.Processes()

	if err != nil {
		panic(fmt.Sprintf("Codecrafters Internal Error - Error fetching process tree: %s", err))
	}

	for _, p := range processes {
		parent, err := p.Ppid()
		if err != nil {
			continue
		}

		if parent != int32(shellPid) {
			continue
		}

		cmd, err := p.Cmdline()

		if err != nil {
			continue
		}

		if cmd == cmdLine {
			return int(p.Pid)
		}
	}

	return -1
}
