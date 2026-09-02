package caldav

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newURLCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "url",
		Short: "Print the CalDAV home set URL",
		Long:  "Print the CalDAV home set URL for the current user, suitable for pasting into Apple Calendar, iOS, or Thunderbird.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			if config.Email == "" {
				return fmt.Errorf("no email configured; run `agenda whoami` or login first")
			}

			caldavURL := fmt.Sprintf("%s/dav/%s/calendars", config.ServerURL, config.Email)
			ui.PrintSuccess("CalDAV home set URL:")
			fmt.Println(caldavURL)

			return nil
		},
	}
}
