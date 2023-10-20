package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env    string   `yaml:"env"`
	Server Server   `yaml:"server"`
	DB     DBConfig `yaml:"database"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DBConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Address  string `yaml:"address"`
	SSLMode  string `yaml:"ssl_mode"`
}

func New(path string) (*Config, error) {
	const op = "config.New"
	err := validConfigPath(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s: file is not open. %w", op, err)
	}
	defer func() { _ = file.Close() }()

	var cfg *Config
	configDecoder := yaml.NewDecoder(file)

	if err := configDecoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("%s: decode fault. %w", op, err)
	}

	return cfg, nil
}

func validConfigPath(path string) error {
	const op = "config.ValidConfigPath"

	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s: it is directory. You need file", op)
	}

	return nil
}

func (cfg *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
}
