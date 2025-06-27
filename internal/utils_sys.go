package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getCommandName(pid int) string {
	commPath := fmt.Sprintf("/proc/%d/comm", pid)
	if commData, err := os.ReadFile(commPath); err == nil {
		return strings.TrimSpace(string(commData))
	}
	return "unknown"
}

func getParentPid(pid int) int {
	if statData, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid)); err == nil {
		fields := strings.Fields(string(statData))
		if len(fields) > 3 {
			if parentPid, err := strconv.Atoi(fields[3]); err == nil {
				return parentPid
			}
		}
	}
	return 0
}

func isTestingTesterUsingBusyboxOnAlpine() bool {

	_, err := os.Stat("/etc/alpine-release")
	isOnAlpine := err == nil

	if !isOnAlpine {
		return false
	}

	pid := os.Getppid()
	fmt.Println("⛳ ppid:", pid)
	current_name := getCommandName(pid)

	parent_pid := getParentPid(pid)
	parent_name := getCommandName(parent_pid)

	great_parent_pid := getParentPid(parent_pid)
	great_parent_name := getCommandName(great_parent_pid)

	fmt.Println("⛳ Current:", current_name)
	fmt.Println("⛳ Parent:", parent_name)
	fmt.Println("⛳ Great-Parent:", great_parent_name)

	isTestingTesterUsingBusybox := current_name == "go" && parent_name == "make" && great_parent_name == "ash"

	return true
	return isTestingTesterUsingBusybox
}
