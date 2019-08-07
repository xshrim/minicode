package sdk

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
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

func Hash(args ...string) string {
	hash := sha256.New()
	hash.Write([]byte(strings.Join(args, ":")))
	md := hash.Sum(nil)
	return hex.EncodeToString(md)
}
