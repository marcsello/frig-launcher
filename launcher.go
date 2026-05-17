package main

import (
	"log"
	"os"
	"os/exec"
)

func Launch(command []string) error {
	log.Println("launching:", command)
	cwd, err := os.Getwd() // reuse CWD for now...
	if err != nil {
		return err
	}

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = cwd

	err = cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Process.Release()
	if err != nil {
		return err
	}

	return nil
}
