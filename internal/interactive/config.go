package interactive

import (
	"fmt"
	"github.com/byawitz/ggh/internal/config"
	"github.com/byawitz/ggh/internal/ssh"
	"github.com/charmbracelet/bubbles/table"
	"os"
)

func InteractiveConfig(value string) []string {
	list, err := config.ParseWithSearch(value, config.GetConfigFile())
	if err != nil || len(list) == 0 {
		fmt.Println("No config found.")
		os.Exit(0)
	}

	var rows []table.Row
	for _, c := range list {
		rows = append(rows, table.Row{
			c.Name,
			c.Host,
			c.Port,
			c.User,
			c.Key,
		})
	}
	c := Select(rows, SelectConfig)
	return ssh.GenerateCommandArgs(c)
}
