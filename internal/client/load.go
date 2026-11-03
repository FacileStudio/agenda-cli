package client

import (
	"net/url"
	"os"
	"strings"
)

// DefaultServerURL is the default Agenda API endpoint.
const DefaultServerURL = "https://agenda.facile.studio/api"

// NormalizeServerURL normalizes server URLs, ensuring default Agenda host points to /api.
func NormalizeServerURL(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimRight(raw, "/")
	if raw == "" {
		return DefaultServerURL
	}
	if raw == "https://agenda.facile.studio" || raw == "http://agenda.facile.studio" {
		return DefaultServerURL
	}
	return raw
}

// NormalizeLoginURL turns a server address typed on the command line into an API URL, adding /api to a bare host.
func NormalizeLoginURL(raw string) string {
	raw = strings.TrimRight(strings.TrimSpace(raw), "/")
	if u, err := url.Parse(raw); err == nil && u.Host != "" && u.Path == "" {
		raw += "/api"
	}
	return NormalizeServerURL(raw)
}

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

	if config.ServerURL == "" {
		config.ServerURL = DefaultServerURL
	} else {
		config.ServerURL = NormalizeServerURL(config.ServerURL)
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
