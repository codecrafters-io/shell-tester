package shell_executable

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

// GetAllDescendentsPids returns the PIDs of every descendant of the shell process (recursive).
// This is useful because the user code may run under a launcher (e.g. uv, another interpreter);
// the real shell is often a child, and background jobs may be deeper in the tree.
func (b *ShellExecutable) GetAllDescendentsPids() []int {
	shellPid := b.GetPid()

	rootShellProcess, err := process.NewProcess(int32(shellPid))

	// This will never trigger: b.GetPid() ensures this
	if err != nil {
		panic(fmt.Sprintf("Codecrafters Internal Error - Could not create Process struct from PID %d: %s", shellPid, err))
	}

	var allDescendantPids []int
	visitedDescendantPids := make(map[int32]struct{})

	var appendDescendantsOf func(parentProcess *process.Process)

	appendDescendantsOf = func(parentProcess *process.Process) {
		directChildren, err := parentProcess.Children()
		if err != nil {
			return
		}

		for _, childProcess := range directChildren {
			childPid := childProcess.Pid

			if _, alreadyRecorded := visitedDescendantPids[childPid]; alreadyRecorded {
				continue
			}

			visitedDescendantPids[childPid] = struct{}{}
			allDescendantPids = append(allDescendantPids, int(childPid))
			appendDescendantsOf(childProcess)
		}
	}

	appendDescendantsOf(rootShellProcess)
	return allDescendantPids
}
