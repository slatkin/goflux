package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type ThemeConfig struct {
	UnreadColor string `toml:"unread_color"`
	ReadColor   string `toml:"read_color"`
}

func DefaultThemeConfig() ThemeConfig {
	return ThemeConfig{
		UnreadColor: "reset", // Placeholder for actual color codes or hex
		ReadColor:   "gray",
	}
}

type Config struct {
	ApiKey            string      `toml:"api_key"`
	ServerUrl         string      `toml:"server_url"`
	AllowInvalidCerts bool        `toml:"allow_invalid_certs"`
	Theme             ThemeConfig `toml:"theme"`
}

func DefaultConfig() Config {
	return Config{
		ApiKey:            "FIXME",
		ServerUrl:         "FIXME",
		AllowInvalidCerts: false,
		Theme:             DefaultThemeConfig(),
	}
}

func GetConfigFilepath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not find user config dir: %w", err)
	}
	// Matching Rust: directories::ProjectDirs::from("com", "spencerwi", "cliflux")
	// usually maps to ~/.config/cliflux on Linux
	path := filepath.Join(configDir, "cliflux", "config.toml")
	return path, nil
}

func Init() (string, error) {
	path, err := GetConfigFilepath()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("configuration file already exists at %s", path)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(DefaultConfig()); err != nil {
		return "", err
	}

	return path, nil
}

func Load() (Config, error) {
	path, err := GetConfigFilepath()
	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config file not found at %s. run with --init to create one", path)
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %w", err)
	}

	// Validate/Clean URL
	cfg.ServerUrl = strings.TrimSpace(cfg.ServerUrl)
	if cfg.ServerUrl == "" {
		return Config{}, errors.New("server_url cannot be empty")
	}
	if strings.HasSuffix(cfg.ServerUrl, "/") {
		cfg.ServerUrl = strings.TrimSuffix(cfg.ServerUrl, "/")
	}

	if cfg.Theme.UnreadColor == "" {
		cfg.Theme.UnreadColor = DefaultThemeConfig().UnreadColor
	}
	if cfg.Theme.ReadColor == "" {
		cfg.Theme.ReadColor = DefaultThemeConfig().ReadColor
	}

	return cfg, nil
}
