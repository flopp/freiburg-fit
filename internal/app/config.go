package app

import (
	"encoding/json"
	"os"
)

type Config struct {
	IsRemoteTarget bool
	OutputDir      string
	UmamiId        string
}

func LoadConfig(jsonFileName string) (Config, error) {
	var config Config
	file, err := os.Open(jsonFileName)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
