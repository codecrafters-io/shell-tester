//go:build !linux

package shell_executable

import (
	"os/exec"
)

// newMemoryMonitor returns a no-op monitor on non-Linux platforms
func newMemoryMonitor(memoryLimitBytes int64) *memoryMonitor {
	return &memoryMonitor{}
}

// memoryMonitor is a no-op on non-Linux platforms
type memoryMonitor struct{}

func (m *memoryMonitor) start(pid int)        {}
func (m *memoryMonitor) stop()                {}
func (m *memoryMonitor) wasOOMKilled() bool   { return false }

// configureProcAttr is a no-op on non-Linux platforms
func configureProcAttr(cmd *exec.Cmd) {}
