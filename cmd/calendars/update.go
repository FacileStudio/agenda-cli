package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var (
		name       string
		color      string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a calendar",
		Long:  "Update a calendar's name or color.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(args[0], name, color, jsonOutput)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New calendar name")
	cmd.Flags().StringVar(&color, "color", "", "New calendar color")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runUpdate(id, name, color string, jsonOutput bool) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	body := make(map[string]any)
	if name != "" {
		body["name"] = name
	}
	if color != "" {
		body["color"] = color
	}

	if len(body) == 0 {
		return fmt.Errorf("nothing to update; use --name or --color")
	}

	resp, err := api.Put("/calendars/"+id, body)
	if err != nil {
		return err
	}

	var calendar map[string]any
	if err := api.DecodeJSON(resp, &calendar); err != nil {
		return err
	}

	if jsonOutput {
		return ui.PrintJSON(calendar)
	}

	ui.PrintSuccess("Calendar updated")
	ui.PrintKeyValue("ID", fmt.Sprintf("%v", calendar["id"]))
	ui.PrintKeyValue("Name", fmt.Sprintf("%v", calendar["name"]))
	ui.PrintKeyValue("Color", fmt.Sprintf("%v", calendar["color"]))

	return nil
}
