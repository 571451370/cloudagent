package main

import (
	"os"
	"path/filepath"

	"github.com/vtolstov/go-ps"
)

func getPids(name string, filter bool) []int {
	name = filepath.Base(name)
	pids := []int{}
	procs, err := ps.FindProcessByExecutable(name)

	if err != nil {
		return pids
	}

	ownpid := os.Getpid()
Check:
	for _, proc := range procs {
		if filter {
			for _, pid := range proc.CPids() {
				if pid == ownpid {
					continue Check
				}
			}
		}
		pids = append(pids, proc.Pid())
		pids = append(pids, proc.CPids()...)
	}

	return pids
}
