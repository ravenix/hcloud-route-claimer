package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Load(path string) (*Config, error) {
	configFile, err := os.Open(path)
	defer configFile.Close()

	if err != nil {
		return nil, err
	}

	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	json.Unmarshal(configBytes, &config)

	return &config, nil
}