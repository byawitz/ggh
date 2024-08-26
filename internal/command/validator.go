package command

import (
	"log"
	"os/exec"
)

func CheckSSH() {
	_, err := exec.LookPath("ssh")
	if err != nil {
		log.Fatal("ssh is not installed")
	}
}
