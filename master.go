package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"
)

func master() error {
	var err error

	master := true
	if os.Getenv("master") == "false" {
		master = false
	}

	if err = os.Chdir("/"); err != nil {
		return err
	}

	unix.Umask(0)

	if master {
		if err = ioutil.WriteFile("/proc/self/oom_score_adj", []byte("-1000"), 0644); err != nil {
			if err = ioutil.WriteFile("/proc/self/oom_adj", []byte("-17"), 0644); err != nil {
				return err
			}
		}

		if err = unix.Setpgid(0, 0); err != nil {
			//l.Info("setpgid: " + err.Error())
		}

		if _, err = unix.Setsid(); err != nil {
			//l.Info("setsid: " + err.Error())
		}

		for _, pid := range getPids(os.Args[0], true) {
			unix.Kill(pid, unix.SIGTERM)
		}
		for _, pid := range getPids(filepath.Base(os.Args[0]), true) {
			unix.Kill(pid, unix.SIGTERM)
		}

	}

	if master {
		stdOut := bytes.NewBuffer(nil)
		stdErr := bytes.NewBuffer(nil)
		for {
			cmd := exec.Command(os.Args[0])
			cmd.Dir = "/"
			cmd.Env = append(cmd.Env, os.Environ()...)
			cmd.Env = append(cmd.Env, "master=false")
			cmd.Stdin = nil
			cmd.Stdout = stdOut
			cmd.Stderr = stdErr
			cmd.ExtraFiles = nil
			cmd.Run()
			if stdErr.Len() > 0 {
				l.Error(string(stdErr.Bytes()))
				stdErr.Reset()
			}
			if stdOut.Len() > 0 {
				l.Info(string(stdOut.Bytes()))
				stdOut.Reset()
			}
			time.Sleep(5 * time.Second)
		}
	} else {
		//		for {
		if err = slave(); err != nil {
			return err
		}
		//			time.Sleep(5 * time.Second)
		//		}
	}
	return nil
}
