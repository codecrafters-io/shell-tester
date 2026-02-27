package shell_executable

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

func (b *ShellExecutable) GetAllChildrenPids() []int {
	shellPid := b.GetPid()

	proc, err := process.NewProcess(int32(shellPid))

	// This will never trigger: b.GetPid() ensures this
	if err != nil {
		panic(fmt.Sprintf("Codecrafters Internal Error - Could not create Process struct from PID %d: %s", shellPid, err))
	}

	children, err := proc.Children()

	if err != nil {
		// Could not get children, return empty slice
		return []int{}
	}

	var pids []int

	for _, child := range children {
		pids = append(pids, int(child.Pid))
	}

	return pids
}
