package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"util-pipe/internal/utils"
	"util-pipe/internal/xj"

	"gopkg.in/yaml.v3"
)

var config Config

func init() {}

type Config struct {
	ScriptsPath string `json:"scripts_path,omitempty" yaml:"scripts_path,omitempty"`
	WikiPath    string `json:"wiki_path,omitempty" yaml:"wiki_path,omitempty"`
	Server      struct {
		Host string `json:"host,omitempty" yaml:"host,omitempty"`
		SSL  struct {
			CRT string `json:"crt,omitempty" yaml:"crt,omitempty"`
			Key string `json:"key,omitempty" yaml:"key,omitempty"`
		} `json:"ssl,omitempty" yaml:"ssl,omitempty"`
	} `json:"server,omitempty" yaml:"server,omitempty"`
	BasicAuth struct {
		Login string `json:"login,omitempty" yaml:"login,omitempty"`
		Pass  string `json:"pass,omitempty" yaml:"pass,omitempty"`
	} `json:"basic_auth,omitempty" yaml:"basic_auth,omitempty"`
	LDAP struct {
		URL string `json:"url,omitempty" yaml:"url,omitempty"`
		DN  string `json:"dn,omitempty" yaml:"dn,omitempty"`
	} `json:"ldap,omitempty" yaml:"ldap,omitempty"`
	wiki struct{ ListLinks []string }
}

func (c *Config) Init(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return c.Unmarshal(filepath.Ext(path), content)
}

func (c *Config) Marshal(format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "j", "json", ".json":
		return xj.Parser.MarshalIndent(c, "", "    ")
	case "x", "xml", ".xml":
		return xml.MarshalIndent(c, "", "    ")
	case "y", "yaml", ".yaml":
		return yaml.Marshal(c)
	}
	return nil, errors.New("unknown configuration format")
}

func (c *Config) Unmarshal(format string, data []byte) error {
	switch strings.ToLower(format) {
	case "j", "json", ".json":
		return xj.Parser.Unmarshal(data, c)
	case "x", "xml", ".xml":
		return xml.Unmarshal(data, c)
	case "y", "yaml", ".yaml":
		return yaml.Unmarshal(data, c)
	}
	return errors.New("unknown configuration format")
}

func (c *Config) Example(format string) string {
	c.ScriptsPath = "/root/scripts"
	c.WikiPath = "/root/wiki"
	c.Server.Host = "127.0.0.1:8090"
	c.Server.SSL.CRT = "server.crt"
	c.Server.SSL.Key = "server.key"
	c.BasicAuth.Login = "admin"
	c.BasicAuth.Pass = "qwerty"
	c.LDAP.URL = "ldap://example.com:389"
	c.LDAP.DN = "ou=users,dc=example,dc=com"
	if runtime.GOOS == "windows" {
		c.ScriptsPath = "c:\\scripts"
		c.WikiPath = "c:\\wiki"
	}
	res, err := c.Marshal(format)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(res)
}

func (c *Config) Check() error {
	if len(c.ScriptsPath) > 0 && !utils.IsDir(c.ScriptsPath) {
		return errors.New("scripts path not found or not dir")
	}
	if len(c.WikiPath) > 0 && !utils.IsDir(c.WikiPath) {
		return errors.New("wiki path not found or not dir")
	}
	if len(c.BasicAuth.Login)+len(c.BasicAuth.Pass) > 0 {
		if len(c.BasicAuth.Login) < 1 {
			return errors.New("basic auth login can't be empty")
		}
		if len(c.BasicAuth.Pass) < 1 {
			return errors.New("basic auth pass can't be empty")
		}
	}
	if len(c.LDAP.URL)+len(c.LDAP.DN) > 0 {
		if len(c.LDAP.URL) < 1 {
			return errors.New("ldap url can't be empty")
		}
		if len(c.LDAP.DN) < 1 {
			return errors.New("ldap dn can't be empty")
		}
	}
	if len(c.Server.SSL.CRT)+len(c.Server.SSL.Key) > 0 {
		if !utils.IsFile(c.Server.SSL.CRT) {
			return errors.New("ssl .crt not found")
		}
		if !utils.IsFile(c.Server.SSL.Key) {
			return errors.New("ssl .key not found")
		}
	}
	return nil
}
