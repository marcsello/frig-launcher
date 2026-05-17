package launcher

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func startCMD(command []string) (*exec.Cmd, error) {
	log.Println("launching subprocess:", command)
	cwd, err := os.Getwd() // reuse CWD for now...
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = cwd

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func startExec(command []string) error {
	log.Println("launching by exec:", command)
	return syscall.Exec(command[0], command[1:], os.Environ())
}
