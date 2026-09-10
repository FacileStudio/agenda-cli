package calendars

import (
	"github.com/FacileStudio/agenda-cli/internal/client"
)

// loadAPI loads the persisted config and returns an authenticated client.
func loadAPI() (*client.Client, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return client.NewClient(config), nil
}

// loadConfig reads the persisted client configuration.
func loadConfig() (*client.Config, error) {
	configPath := client.DefaultConfigPath()
	config := &client.Config{}
	if err := config.Read(configPath); err != nil {
		return nil, err
	}
	return config, nil
}
