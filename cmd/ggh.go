package cmd

import (
	"fmt"
	"github.com/byawitz/ggh/internal/command"
	"github.com/byawitz/ggh/internal/config"
	"github.com/byawitz/ggh/internal/history"
	"github.com/byawitz/ggh/internal/interactive"
	"github.com/byawitz/ggh/internal/ssh"
	"os"
)

func Main() {
	command.CheckSSH()

	args := os.Args[1:]

	fmt.Println("\033[2mIn memory of Binyamin Yawitz (1990–2025), creator of GGH \033[31m❤️\033[0m\033[2m\033[0m")

	action, value := command.Which()
	switch action {
	case command.InteractiveHistory:
		args = interactive.History()
	case command.InteractiveConfig:
		args = interactive.Config("")
	case command.InteractiveConfigWithSearch:
		args = interactive.Config(value)
	case command.ListHistory:
		history.Print()
		return
	case command.ListConfig:
		config.Print()
		return
	default:
		history.AddHistoryFromArgs(args)
	}
	ssh.Run(args)
}
