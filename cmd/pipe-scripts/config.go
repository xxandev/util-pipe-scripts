package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"runtime"
	"util-pipe-scripts/internal/utils"

	"gopkg.in/yaml.v2"
)

type FormatConfig int8

const (
	JSON FormatConfig = iota
	XML
	YAML
)

type Config struct {
	Host string `json:"server_host,omitempty" xml:"server_host,omitempty" yaml:"server_host,omitempty"`
	Path string `json:"scripts_path,omitempty" xml:"scripts_path,omitempty" yaml:"scripts_path,omitempty"`
}

func (c *Config) Init() (bool, error) {
	switch {
	case utils.IsFile("pipe-scripts.json"):
		content, err := ioutil.ReadFile("pipe-scripts.json")
		if err != nil {
			return true, err
		}
		if err := c.Unmarshal(JSON, content); err != nil {
			return true, err
		}
	case utils.IsFile("pipe-scripts.xml"):
		content, err := ioutil.ReadFile("pipe-scripts.xml")
		if err != nil {
			return true, err
		}
		if err := c.Unmarshal(XML, content); err != nil {
			return true, err
		}
	case utils.IsFile("pipe-scripts.yaml"):
		content, err := ioutil.ReadFile("pipe-scripts.yaml")
		if err != nil {
			return true, err
		}
		if err := c.Unmarshal(YAML, content); err != nil {
			return true, err
		}
	}
	return false, errors.New("config file not found")
}

func (c *Config) Marshal(format FormatConfig) ([]byte, error) {
	switch format {
	case JSON:
		return json.MarshalIndent(c, "", "\t")
	case XML:
		return xml.MarshalIndent(c, "", "\t")
	case YAML:
		return yaml.Marshal(c)
	}
	return nil, errors.New("unknown configuration format")
}

func (c *Config) Unmarshal(format FormatConfig, data []byte) error {
	switch format {
	case JSON:
		return json.Unmarshal(data, c)
	case XML:
		return xml.Unmarshal(data, c)
	case YAML:
		return yaml.Unmarshal(data, c)
	}
	return errors.New("unknown configuration format")
}

func (c *Config) Example() {
	c.Host = "127.0.0.1:8090"
	if runtime.GOOS == "windows" {
		c.Path = "c:\\scripts"
	} else {
		c.Path = "/root/scripts"
	}
}

func (c *Config) Check() error {
	return nil
}
