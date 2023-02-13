package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed config.yaml
var Default []byte

var defaultInstance = func() *Config {
	var cfg Config
	err := yaml.Unmarshal(Default, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}()

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
	switch {
	case os.IsNotExist(err):
		dir := filepath.Dir(path)
		_, err = os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					return nil, fmt.Errorf("failed to create config directory: %v", err)
				}
			} else {
				return nil, fmt.Errorf("config directory error: %v", err)
			}
		}

		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		_, err = file.Write(Default)
		if err != nil {
			return nil, fmt.Errorf("failed to write default config: %v", err)
		}
		logrus.Warn("cannot find config file, write default content")
		return defaultInstance, nil

	case err == nil:
		if stat.IsDir() {
			return nil, fmt.Errorf("config file %s is a directory", path)
		}

	default:
		return nil, fmt.Errorf("failed to read config file: %v", err)
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
