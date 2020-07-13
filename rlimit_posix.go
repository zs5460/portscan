// +build !windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

func checkLimit() {
	const min = 10240

	rlimit := &syscall.Rlimit{}
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, rlimit)
	if err == nil && rlimit.Cur < min {
		fmt.Printf("WARNING: File descriptor limit %d is too low. "+
			"At least %d is recommended. Fix with `ulimit -n %d`.\n", rlimit.Cur, min, min)
	}
	os.Exit(-1)
}
