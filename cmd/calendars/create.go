package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var color string

	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new calendar",
		Long:  "Create a new calendar with the given name and optional color.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			api := client.NewClient(config)
			body := map[string]interface{}{"name": name}
			if color != "" {
				body["color"] = color
			}

			resp, err := api.Post("/calendars", body)
			if err != nil {
				return err
			}

			var calendar map[string]interface{}
			if err := api.DecodeJSON(resp, &calendar); err != nil {
				return err
			}

			ui.PrintSuccess("Calendar created")
			ui.PrintKeyValue("ID", fmt.Sprintf("%v", calendar["id"]))
			ui.PrintKeyValue("Name", fmt.Sprintf("%v", calendar["name"]))
			ui.PrintKeyValue("Color", fmt.Sprintf("%v", calendar["color"]))

			return nil
		},
	}

	cmd.Flags().StringVar(&color, "color", "", "Calendar color (hex or name)")

	return cmd
}
