package main

import (
	"errors"
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddress string      `yaml:"listen_address"`
	UseHTTPS      bool        `yaml:"use_https"`
	KeyFile       string      `yaml:"keyfile"`
	CertFile      string      `yaml:"certfile"`
	Hostname      string      `yaml:"hostname"`
	MySQL         mysqlConfig `yaml:"mysql"`
	DebugMode     bool        `yaml:"debug_mode"`
}

type mysqlConfig struct {
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func NewDefaultConfig() *Config {
	return &Config{
		ListenAddress: ":80",
		KeyFile:       "",
		CertFile:      "",
		UseHTTPS:      false,
		Hostname:      "localhost",
		MySQL: mysqlConfig{
			Database: "fakebook",
			User:     "fakebook",
			Password: "",
		},
		DebugMode: true,
	}
}

func ReadConfigFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.UseHTTPS && (config.KeyFile == "" || config.CertFile == "") {
		return nil, errors.New(
			"neither <keyfile> or <certfile> is specified with <use_https> flag set")
	}

	return &config, nil
}

func (c *Config) GetPortString() string {
	_, port, _ := net.SplitHostPort(c.ListenAddress)
	return port
}

func (c *Config) BasicURL() string {
	scheme := "http"
	if c.UseHTTPS {
		scheme = "https"
	}

	port := c.GetPortString()

	optionalPort := ""
	if !isDefaultProtocolPort(scheme, port) {
		optionalPort = ":" + port
	}

	return fmt.Sprintf("%s://%s%s", scheme, c.Hostname, optionalPort)
}

func isDefaultProtocolPort(proto, port string) bool {
	return (proto == "http" && port == "80") ||
		(proto == "https" && port == "443")
}
