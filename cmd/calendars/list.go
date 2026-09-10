package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all calendars",
		Long:  "List all calendars visible to the authenticated user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			api := client.NewClient(config)
			resp, err := api.Get("/calendars")
			if err != nil {
				return err
			}

			var calendars []map[string]any
			if err := api.DecodeJSON(resp, &calendars); err != nil {
				return err
			}

			ui.PrintSuccess(fmt.Sprintf("Found %d calendars", len(calendars)))
			for _, cal := range calendars {
				ui.PrintKeyValue("ID", fmt.Sprintf("%v", cal["id"]))
				ui.PrintKeyValue("Name", fmt.Sprintf("%v", cal["name"]))
				ui.PrintKeyValue("Color", fmt.Sprintf("%v", cal["color"]))
				fmt.Println()
			}

			return nil
		},
	}
}
