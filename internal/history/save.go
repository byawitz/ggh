package history

import (
	"encoding/json"
	"fmt"
	"github.com/byawitz/ggh/internal/config"
	"github.com/charmbracelet/bubbles/table"
	"os"
	"strings"
	"time"
)

func AddHistoryFromArgs(args []string) {
	if len(args) == 1 && !strings.Contains(args[0], "@") {
		localConfig, err := config.GetConfig(args[0])
		if err != nil || localConfig.Name == "" {
			return
		}

		AddHistory(localConfig)
		return
	}

	generatedConfig := config.SSHConfig{}

	skipNext := false
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}

		switch {
		case strings.HasPrefix(arg, "-p"):
			if arg == "-p" {
				generatedConfig.Port = args[i+1]
				skipNext = true
			} else {
				generatedConfig.Port = args[i][2:]
			}
		case arg == "-i":
			generatedConfig.Key = args[i+1]
			skipNext = true
		case strings.Contains(arg, "@"):
			values := strings.Split(arg, "@")
			generatedConfig.User = values[0]
			generatedConfig.Host = values[1]
		}
	}
	AddHistory(generatedConfig)
}

func AddHistory(c config.SSHConfig) {
	if c.Host == "" {
		return
	}

	list, err := Fetch(getFile())

	if err != nil {
		fmt.Println("error getting ggh file")
		return
	}

	err = saveFile(SSHHistory{Connection: c, Date: time.Now()}, list)
	if err != nil {
		fmt.Println("error saving ggh file")
		return
	}
}

func RemoveByIP(row table.Row) {
	list, err := Fetch(getFile())

	if err != nil {
		fmt.Println("error getting ggh file")
		return
	}

	ip := row[1]

	saving := make([]SSHHistory, 0, len(list)-1)

	for _, item := range list {
		if item.Connection.Host == ip {
			continue
		}

		saving = append(saving, item)
	}

	err = saveFile(SSHHistory{}, saving)
	if err != nil {
		panic("error saving ggh file")
	}

}

func saveFile(n SSHHistory, l []SSHHistory) error {
	file := getFileLocation()
	fileContent := stringify(n, l)

	err := os.WriteFile(file, []byte(fileContent), 0644)

	return err
}

func stringify(n SSHHistory, l []SSHHistory) string {
	history := make([]SSHHistory, 0)

	if n.Connection.Host != "" {
		history = append(history, n)
	}

	for _, sshHistory := range l {
		sshHistory.Connection.CleanName()
		if sshHistory.Connection.UniqueKey() != n.Connection.UniqueKey() {
			history = append(history, sshHistory)
		}
	}

	content, err := json.Marshal(history)

	if err != nil {
		return ""
	}

	return string(content)
}
