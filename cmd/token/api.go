package token

import (
	"github.com/FacileStudio/agenda-cli/internal/client"
)

func loadAPI() (*client.Client, error) {
	return client.LoadAPI()
}
