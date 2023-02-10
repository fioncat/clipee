package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	defaultListenAddr = ":9988"
)

type Config struct {
	Listen  string   `yaml:"listen"`
	Remotes []string `yaml:"remotes"`

	Share string `yaml:"share"`
}

var (
	instance *Config
	initOnce sync.Once
)

func doInit() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := os.Getenv("CLIPEE_CONFIG")
	if path == "" {
		path = filepath.Join(homeDir, ".config", "clipee", "config.yaml")
	}

	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("please create config file %s first", path)
		}
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	if stat.IsDir() {
		return nil, fmt.Errorf("config file %s is a directory", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read and parse config file: %v", err)
	}

	if cfg.Listen == "" {
		logrus.Warnf("use default listen address %s", defaultListenAddr)
		cfg.Listen = defaultListenAddr
	}

	if cfg.Share == "" {
		share := filepath.Join(homeDir, "Desktop", "share")
		logrus.Warnf("use default share path %s", share)
		cfg.Share = share
	}
	cfg.Share = os.ExpandEnv(cfg.Share)

	return &cfg, nil
}

func Init() error {
	var err error
	initOnce.Do(func() {
		instance, err = doInit()
	})
	return err
}

func Get() *Config {
	if instance == nil {
		panic("please call Init first")
	}
	return instance
}
