package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show calendar details",
		Long:  "Show details for a specific calendar by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			api := client.NewClient(config)
			resp, err := api.Get("/calendars/" + id)
			if err != nil {
				return err
			}

			var calendar map[string]any
			if err := api.DecodeJSON(resp, &calendar); err != nil {
				return err
			}

			ui.PrintKeyValue("ID", fmt.Sprintf("%v", calendar["id"]))
			ui.PrintKeyValue("Name", fmt.Sprintf("%v", calendar["name"]))
			ui.PrintKeyValue("Color", fmt.Sprintf("%v", calendar["color"]))
			ui.PrintKeyValue("Owner", fmt.Sprintf("%v", calendar["owner_id"]))

			return nil
		},
	}
}
