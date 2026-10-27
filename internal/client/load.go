package client

import (
	"os"
)

// Load loads configuration from disk and environment variables.
func Load() (*Config, error) {
	config := &Config{}
	path := DefaultConfigPath()
	if path != "" {
		if err := config.Read(path); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	if envURL := os.Getenv("AGENDA_SERVER_URL"); envURL != "" {
		config.ServerURL = envURL
	}
	if envToken := os.Getenv("AGENDA_TOKEN"); envToken != "" {
		config.Token = envToken
	}

	return config, nil
}

// LoadAPI loads configuration and returns an authenticated client.
func LoadAPI() (*Client, error) {
	config, err := Load()
	if err != nil {
		return nil, err
	}
	if err := config.RequireAuth(); err != nil {
		return nil, err
	}
	return NewClient(config), nil
}
