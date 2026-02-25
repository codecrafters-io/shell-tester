//go:build linux

package shell_executable

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// memoryMonitor monitors process memory usage via /proc and kills if limit exceeded
type memoryMonitor struct {
	pid       int
	limit     int64
	oomKilled atomic.Bool
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

// newMemoryMonitor creates a new memory monitor with the specified limit.
// Call startMonitoring() after the process has started to begin monitoring.
func newMemoryMonitor(memoryLimitBytes int64) *memoryMonitor {
	return &memoryMonitor{
		limit: memoryLimitBytes,
	}
}

// start begins polling /proc for RSS usage of the given process.
// Must be called after the process has started.
func (m *memoryMonitor) start(pid int) {
	if m.limit <= 0 {
		return
	}

	m.pid = pid
	m.stopChan = make(chan struct{})
	m.wg.Add(1)
	go m.monitor()
}

// monitor polls the process tree's RSS and kills if limit is exceeded
func (m *memoryMonitor) monitor() {
	defer m.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			rss, err := getProcessTreeRSS(m.pid)
			if err != nil {
				// Process likely exited, stop monitoring
				return
			}

			if rss > m.limit {
				m.oomKilled.Store(true)
				// Kill the process group to ensure all children are terminated
				// We are replying on `cmd.SysProcAttr.Setsid = true` in creack/pty's StartWithSize to avoid process group conflicts.
				// https://github.com/creack/pty/blob/master/start.go
				syscall.Kill(-m.pid, syscall.SIGKILL)
				syscall.Kill(m.pid, syscall.SIGKILL)
				return
			}
		}
	}
}

// wasOOMKilled returns true if the process was killed due to exceeding memory limit
func (m *memoryMonitor) wasOOMKilled() bool {
	return m.oomKilled.Load()
}

// stop stops the memory monitor
func (m *memoryMonitor) stop() {
	if m.stopChan != nil {
		close(m.stopChan)
		m.wg.Wait()
		m.stopChan = nil
	}
}

// getProcessTreeRSS returns the total RSS (in bytes) of a process and all its descendants
func getProcessTreeRSS(pid int) (int64, error) {
	visited := make(map[int]bool)
	return getProcessTreeRSSRecursive(pid, visited)
}

func getProcessTreeRSSRecursive(pid int, visited map[int]bool) (int64, error) {
	if visited[pid] {
		return 0, nil
	}
	visited[pid] = true

	// Get RSS for this process
	rss, err := getProcessRSS(pid)
	if err != nil {
		return 0, err
	}

	// Find and sum RSS of all children
	children, err := getChildPIDs(pid)
	if err != nil {
		// If we can't read children, just return this process's RSS
		return rss, nil
	}

	for _, childPID := range children {
		childRSS, err := getProcessTreeRSSRecursive(childPID, visited)
		if err != nil {
			// Child may have exited, continue with others
			continue
		}
		rss += childRSS
	}

	return rss, nil
}

// getProcessRSS reads RSS from /proc/<pid>/statm and returns bytes
func getProcessRSS(pid int) (int64, error) {
	statmPath := fmt.Sprintf("/proc/%d/statm", pid)
	data, err := os.ReadFile(statmPath)
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 2 {
		return 0, fmt.Errorf("unexpected statm format")
	}

	// Second field is RSS in pages
	pages, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return pages * int64(os.Getpagesize()), nil
}

// getChildPIDs returns the PIDs of all direct children of the given process
func getChildPIDs(pid int) ([]int, error) {
	var children []int

	// Read from /proc/<pid>/task/*/children for all threads
	taskDir := fmt.Sprintf("/proc/%d/task", pid)
	tasks, err := os.ReadDir(taskDir)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		childrenPath := filepath.Join(taskDir, task.Name(), "children")
		data, err := os.ReadFile(childrenPath)
		if err != nil {
			continue
		}

		for _, pidStr := range strings.Fields(string(data)) {
			childPID, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}
			children = append(children, childPID)
		}
	}

	return children, nil
}
