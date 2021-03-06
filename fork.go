// +build !windows

package main

import (
	"os"
	"os/exec"
	"syscall"
)

func fork() error {
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Dir = "/"
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.ExtraFiles = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: false, Setpgid: true}

	return cmd.Start()
}
