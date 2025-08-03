package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	configFileName = "/.gatorconfig.json"
)

type Config struct {
	DBUrl	string	`json:"db_url"`
	CurrentUser	string	`json:"current_user_name"`
}

func (cfg *Config) SetUser(u string) error {
	cfg.CurrentUser = u
	err := write(*cfg)
	if err != nil {
		return fmt.Errorf("Error writing to config file: %v", err)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error finding home directory: %v", err)
		return "", err
	}

	return homeDir + configFileName, nil
}

func write(cfg Config) error {
	configFile, err := getConfigFilePath()
	if err != nil {
		log.Printf("Error finding home directory: %v", err)
		return err
	}

	f, err := os.Create(configFile)
	if err != nil {
		log.Printf("Error opening config file %s: %v", configFile, err)
		return err
	}
	defer f.Close()

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return err
	}

	_, err = f.Write(jsonData)
	if err != nil {
		log.Printf("Error writing JSON: %v", err)
		return err
	}

	return nil
}

func Read() (Config, error) {
	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	f, err := os.Open(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("Error opening config file %s: %v", configFile, err)
	}
	defer f.Close()

	var config Config
	jsonParser := json.NewDecoder(f)
	if err := jsonParser.Decode(&config); err != nil {
		return config, fmt.Errorf("Error decoding json: %v", err)
	}
	
	return config, nil
}