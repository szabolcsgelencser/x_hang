package main

import (
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/bitrise-io/go-utils/v2/log"
)

func main() {
	logger := log.NewLogger()

	if len(os.Args) < 2 {
		logger.TErrorf("[hm] No command provided to run")
		os.Exit(1)
	}

	doneCh := make(chan error, 1)

	waitCmd := exec.Command(os.Args[1], os.Args[2:]...)
	waitCmd.Stdout = os.Stdout
	waitCmd.Stderr = os.Stderr

	runFn := func() {
		logger.TPrintf("$ %v", waitCmd.Args)
		doneCh <- waitCmd.Run()
	}

	timeout := 30 * time.Second
	timer := time.NewTimer(timeout)

	go runFn()

	select {
	case err := <-doneCh:
		{
			if err != nil {
				var exitError *exec.ExitError
				if errors.As(err, &exitError) {
					exitCode := exitError.ProcessState.ExitCode()
					logger.TPrintf("[hm] command exited with exitcode %d", exitCode)

					os.Exit(exitCode)
				}

				logger.TErrorf("[hm] command failed to execute: %s", err)

				os.Exit(11)
			}

			logger.TPrintf("[hm] command finished")
			os.Exit(0)
		}
	case <-timer.C:
		logger.TErrorf("[hm] timeout after %s", timeout)
		// killEnvman := exec.Command("killall", "envman", "-SIGQUIT")
		// err := killEnvman.Run()
		// if err != nil {
		// 	logger.TWarnf("[hm] killall failed; %s", err)
		// }
		if err := syscall.Kill(waitCmd.Process.Pid, syscall.SIGQUIT); err != nil {
			logger.TWarnf("[hm] failed to kill process: %s", err)
		}

		time.Sleep(6)
		if err := syscall.Kill(waitCmd.Process.Pid, syscall.SIGKILL); err != nil {
			logger.TWarnf("[hm] failed to kill process: %s", err)
		}
		os.Exit(12)
	}

	os.Exit(0)
}
