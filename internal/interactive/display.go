package interactive

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/byawitz/ggh/internal/config"
	"github.com/byawitz/ggh/internal/history"
	"github.com/byawitz/ggh/internal/ssh"
	"github.com/charmbracelet/bubbles/table"
)

func Config(value string) []string {
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
	history.AddHistory(c)
	if c.Name == "" {
		return ssh.GenerateCommandArgs(c)
	}
	return []string{c.Name}
}

func History() []string {
	list, err := history.FetchWithDefaultFile()

	if err != nil {
		log.Fatal(err)
	}

	if len(list) == 0 {
		fmt.Println("No history found.")
		os.Exit(0)
	}

	var rows []table.Row
	currentTime := time.Now()
	for _, historyItem := range list {
		rows = append(rows, table.Row{
			historyItem.Connection.Name,
			historyItem.Connection.Host,
			historyItem.Connection.Port,
			historyItem.Connection.User,
			historyItem.Connection.Key,
			fmt.Sprintf("%s", history.ReadableTime(currentTime.Sub(historyItem.Date))),
		})
	}
	c := Select(rows, SelectHistory)
	history.AddHistory(c)
	if c.Name == "" {
		return ssh.GenerateCommandArgs(c)
	}
	return []string{c.Name}
}
