package executor

import (
	"log"
	"os/exec"

	"gitlab.com/MikeTTh/env"
)

var (
	wCMD       *exec.Cmd
	launchMode LaunchMode
	execAtExit []string
)

func init() {
	var ok bool
	launchMode, ok = ModeFromString(env.String("FRIG_LAUNCH_MODE", LaunchDetachedStr))
	if !ok {
		panic("Invalid launch mode")
	}
}

func Launch(command []string) error {
	switch launchMode {
	case LaunchWait, LaunchDetached:
		var err error
		wCMD, err = startCMD(command)
		if err != nil {
			return err
		}
		if launchMode == LaunchDetached {
			err = wCMD.Process.Release()
			if err != nil {
				log.Println("Failed to release process:", err)
				return err
			}
		}
	case LaunchExec:
		// register for launching at exit
		execAtExit = command
	}
	return nil
}

func AtExit() {
	switch launchMode {
	case LaunchWait:
		if wCMD != nil {
			err := wCMD.Wait()
			if err != nil {
				log.Println("Error running command:", err)
			}
			log.Println("Subprocess exited cleanly.")
		}
	case LaunchExec:
		if execAtExit != nil {
			err := startExec(execAtExit)
			if err != nil {
				log.Println("Failed to exec process:", err)
			}
		}
	}
}
