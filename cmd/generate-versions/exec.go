package main

import "os/exec"

func execShell(cmd string) ([]byte, error) {
	return exec.Command("sh", "-c", cmd).Output()
}
