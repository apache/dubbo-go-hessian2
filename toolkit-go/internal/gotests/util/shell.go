package util

import "os/exec"

func ShellRun(cmd string) error {
	c := exec.Command("bash", "-c", cmd)
	return c.Run()
}
