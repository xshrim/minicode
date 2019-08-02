package sdk

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

func LoadConfig(filename string) (*Config, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
