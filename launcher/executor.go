package launcher

import (
	"log"
	"os/exec"
)

var (
	wCMD       *exec.Cmd
	launchMode LaunchMode
	execAtExit []string
)

func Init(mode LaunchMode) {
	if launchMode != launchInvalid {
		panic("already initialized")
	}
	launchMode = mode
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
	case LaunchNone:
		log.Println("not launching:", command)
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
