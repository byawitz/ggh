package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"path/filepath"

	"github.com/byawitz/ggh/internal/theme"
	"github.com/charmbracelet/bubbles/table"
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

	for _, config := range configsStrings {
		lines := strings.Split(config, "\n")

		if strings.Trim(lines[0], " ") == "" {
			continue
		}

		sshConfig := SSHConfig{
			Name: lines[0],
			Port: "",
			User: "",
		}

		for _, line := range lines {
			if len(line) == 0 || line[0] == '#' {
				continue
			}

			line = strings.ReplaceAll(strings.TrimLeft(line, " "), "\t", "")
			lineData := strings.Split(line, " ")
			value := ""
			if len(lineData) > 1 {
				value = lineData[1]
			}
			switch {
			case strings.Contains(line, "Include"):
				result, err := ParseInclude(search, value)
				if err != nil {
					panic(err)
				}
				configs = append(configs, result...)
			case strings.Contains(line, "Host"):
				sshConfig.Host = value
			case strings.Contains(line, "Port"):
				sshConfig.Port = value
			case strings.Contains(line, "User"):
				sshConfig.User = value
			case strings.Contains(line, "IdentityFile"):
				sshConfig.Key = value
			}
		}

		if sshConfig.Host == "" || !strings.Contains(sshConfig.Name, search) {
			continue
		}

		configs = append(configs, sshConfig)

	}

	return configs, nil
}

func ParseInclude(search string, path string) ([]SSHConfig, error) {
	var results = make([]SSHConfig, 0)

	var isAbsolute = path[0] == '/' || path[0] == '~'

	var paths []string
	var err error

	if isAbsolute {
		if path[0] == '~' {
			path = filepath.Join(HomeDir(), path[2:])
		}
	} else {
		path = filepath.Join(GetSshDir(), path)
	}

	paths, err = filepath.Glob(path)

	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		info, err := os.Stat(path)

		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			continue
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		items, err := ParseWithSearch(search, string(fileContent))
		if err != nil {
			return nil, err
		}
		results = append(results, items...)
	}

	return results, nil
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
