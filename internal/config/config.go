package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl				string	`json:"db_url"`
	CurrentUserName 	string	`json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := homePath + "/" + configFileName
	return path, nil
}

func Read() *Config {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil
	}

	return &config
}

func (cfg *Config) SetUser(userName string) {
	cfg.CurrentUserName = userName

	configData, err := json.Marshal(*cfg)
	if err != nil {
		return
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		return
	}

	if err := os.WriteFile(configPath, configData, 0600); err != nil {
		return
	}
}