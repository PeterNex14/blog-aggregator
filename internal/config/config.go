package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl				string		`json:"db_url"`
	CurrentUserName		string		`json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	var config Config

	path, err := getConfigFilePath()
	if err != nil {
		return config, err
	}

	file, err := os.Open(path)
	if err != nil {
		return config, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user

	if err := write(*cfg); err != nil {
		return err
	}
	
	return nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(home, configFileName)
	return path, nil
}