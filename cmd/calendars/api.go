package calendars

import (
	"github.com/FacileStudio/agenda-cli/internal/client"
)

// loadAPI loads the configuration and returns an authenticated client.
func loadAPI() (*client.Client, error) {
	return client.LoadAPI()
}

