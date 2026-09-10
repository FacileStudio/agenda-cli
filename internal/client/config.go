package client

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the client's persisted configuration.
type Config struct {
	ServerURL string `yaml:"server_url"`
	Token     string `yaml:"token"`
	UserID    string `yaml:"user_id,omitempty"`
	Email     string `yaml:"email,omitempty"`
}

func (c *Config) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (c *Config) Read(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, c)
}

// DefaultConfigPath returns the config file path under the user's config dir.
func DefaultConfigPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "agenda", "config.yml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "agenda", "config.yml")
}

func (c *Config) Validate() error {
	if c.ServerURL == "" {
		return fmt.Errorf("no server URL configured; run `agenda config server <url>` or set AGENDA_SERVER_URL")
	}
	return nil
}

func (c *Config) IsAuthenticated() bool {
	return c.Token != ""
}

func (c *Config) AuthHeader() string {
	return "Bearer " + c.Token
}
