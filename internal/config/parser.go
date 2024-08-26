package config

import (
	"fmt"
	"github.com/byawitz/ggh/internal/theme"
	"github.com/charmbracelet/bubbles/table"
	"log"
	"strings"
)

type SSHConfig struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Key  string `json:"key"`
}

func Parse(configFile string) ([]SSHConfig, error) {
	return ParseWithSearch("", configFile)
}

func ParseWithSearch(search string, configFile string) ([]SSHConfig, error) {
	configsStrings := strings.Split(strings.ReplaceAll(configFile, "\r\n", "\n"), "Host ")
	var configs = make([]SSHConfig, 0)

	for _, config := range configsStrings[1:] {
		lines := strings.Split(config, "\n")

		if strings.Trim(lines[0], " ") == "" {
			continue
		}

		if search != "" && !strings.Contains(lines[0], search) {
			continue
		}

		sshConfig := SSHConfig{
			Name: lines[0],
			Port: "",
			User: "",
		}

		for _, line := range lines {
			line = strings.ReplaceAll(strings.TrimLeft(line, " "), "\t", "")
			switch {
			case strings.Contains(line, "Host"):
				sshConfig.Host = strings.Split(line, " ")[1]
			case strings.Contains(line, "Port"):
				sshConfig.Port = strings.Split(line, " ")[1]
			case strings.Contains(line, "User"):
				sshConfig.User = strings.Split(line, " ")[1]
			case strings.Contains(line, "IdentityFile"):
				sshConfig.Key = strings.Split(line, " ")[1]
			}

		}
		configs = append(configs, sshConfig)
	}

	return configs, nil
}

func Print() {
	list, err := Parse(GetConfigFile())

	if err != nil {
		log.Fatal(err)
	}

	if len(list) == 0 {
		fmt.Println("No configs found in ~/.ssh/config.")
		return
	}

	var rows []table.Row
	for _, history := range list {
		rows = append(rows, table.Row{history.Name, history.Host, history.Port, history.User, history.Key})
	}
	fmt.Println(theme.PrintTable(rows, theme.PrintConfig))

}
