package config

import (
	"os"
	"path/filepath"
)

func GetConfigFile() string {
	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	sshConfigDir := filepath.Join(userHomeDir, ".ssh")

	config, err := os.ReadFile(filepath.Join(sshConfigDir, "config"))

	return string(config)
}

func GetConfig(name string) (SSHConfig, error) {
	list, err := ParseWithSearch(name, GetConfigFile())
	if err != nil {
		return SSHConfig{}, err
	}

	for _, sshConfig := range list {
		if sshConfig.Name == name {
			return sshConfig, nil
		}
	}
	return SSHConfig{}, nil
}
