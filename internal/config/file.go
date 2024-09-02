package config

import (
	"os"
	"path/filepath"
)

func GetSshDir() string {
	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return filepath.Join(userHomeDir, ".ssh")
}

func GetConfigFile() string {
	sshConfigDir := GetSshDir()

	config, err := os.ReadFile(filepath.Join(sshConfigDir, "config"))
	if err != nil {
		return ""
	}

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
