package cmd

import (
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
