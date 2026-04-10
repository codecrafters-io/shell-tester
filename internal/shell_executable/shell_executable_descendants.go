package shell_executable

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

// GetChildrenAndGrandChildrenPids returns the PIDs of all the children and grandchildren
// of the shell process. This is useful because the user code may run under a launcher
// (e.g. uv, another interpreter); the real
// shell is often a child, and background jobs are grandchildren of the launcher PID.
func (b *ShellExecutable) GetChildrenAndGrandChildrenPids() []int {
	shellPid := b.GetPid()

	proc, err := process.NewProcess(int32(shellPid))

	// This will never trigger: b.GetPid() ensures this
	if err != nil {
		panic(fmt.Sprintf("Codecrafters Internal Error - Could not create Process struct from PID %d: %s", shellPid, err))
	}

	children, err := proc.Children()
	if err != nil {
		return []int{}
	}

	var descendantPids []int

	for _, child := range children {
		descendantPids = append(descendantPids, int(child.Pid))

		grandchildren, err := child.Children()

		if err != nil {
			continue
		}

		for _, gc := range grandchildren {
			descendantPids = append(descendantPids, int(gc.Pid))
		}
	}

	return descendantPids
}
