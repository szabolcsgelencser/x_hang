package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
)

func RunCommandInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}

func main() {
	numTests := 100
	// outs := make([]string, numTests)
	// var wg sync.WaitGroup

	runFn := func(id int) error {
		runCmd := command.NewWithStandardOuts("hangman", "bitrise", "--debug", "--loglevel=\"\"", "run", "--secret-filtering=true", "check")
		// tmpFile, err := os.CreateTemp(os.TempDir(), "")
		// if err != nil {
		// 	log.Errorf("%s", err)
		// 	return err
		// }
		// defer func() {
		// 	if err := tmpFile.Close(); err != nil {
		// 		log.Warnf("%s", err)
		// 	}
		// }()

		// envmanArgs := []string{"envman", "--loglevel", "debug", "--path", tmpFile.Name(), "init", "--clear"}
		// err = RunCommandInDir("", "hangman", envmanArgs...)

		runCmd.AppendEnvs("BITRISE_ANALYTICS_DISABLED=false")
		err := runCmd.Run()
		// outs[id] = out
		if err != nil {
			log.Errorf("%d exited with: %s", id, err)
			return err
		}

		return nil

		// wg.Done()
	}

	hanged := 0
	for i := 0; i < numTests; i++ {
		// wg.Add(1)
		// time.Sleep(1 * time.Duration(i))
		// go runFn(i)

		if err := runFn(i); err != nil {
			var exitError *exec.ExitError
			if errors.As(err, &exitError) {
				exitCode := exitError.ProcessState.ExitCode()
				if exitCode > 10 {
					hanged++
				}
			}
		}
	}
	// wg.Wait()

	log.Printf("hanged: %d", hanged)
	if hanged > 0 {
		os.Exit(13)
	}

	/*
		for i, out := range outs {
			// if !strings.Contains(out, "[hm] timeout after") {
			log.Donef("----- Out %d ------", i)
			log.Printf(out)
			// }
		}
		for i, out := range outs {
			if strings.Contains(out, "[hm] timeout after") {
				log.Donef("----- Hang %d ------", i)
				log.Printf(out)
			}
		}*/

	os.Exit(0)
}
