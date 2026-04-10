package utils

import "os"

func IsOnAlpine() bool {
	_, err := os.Stat("/etc/alpine-release")
	return err == nil
}
