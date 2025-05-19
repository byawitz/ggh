package ssh

import (
	"fmt"
	"github.com/byawitz/ggh/internal/config"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func GenerateCommandArgs(c config.SSHConfig) []string {
	key, port := "", ""
	user := "root"

	if c.User != "" {
		user = c.User
	}

	if c.Key != "" {
		key = "-i " + c.Key
	}

	if c.Port != "" {
		port = "-p " + c.Port
	}
	return strings.Split(fmt.Sprintf("%s %s@%s %s", key, user, c.Host, port), " ")
}

func Run(args []string) {
	args = slices.DeleteFunc(args, func(s string) bool { return s == "" })

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
