package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ConfigParser struct {
	Path string
}

func (c *ConfigParser) Load() (*AppConfig, error) {
	file, err := os.Open(c.Path)
	defer file.Close()

	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var appConfig AppConfig
	err = json.Unmarshal(raw, &appConfig)
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}

func NewConfigParser(path string) *ConfigParser {
	config := &ConfigParser{
		Path: path,
	}
	return config
}
