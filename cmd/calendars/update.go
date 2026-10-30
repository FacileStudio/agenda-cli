package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

var carriedFields = []string{"name", "color", "description", "echo_url"}

func newUpdateCmd() *cobra.Command {
	var (
		fields     calendarFields
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a calendar",
		Long: "Update a calendar. Fields you do not pass keep their stored value; " +
			"pass an empty string to clear color, description or echo URL.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cmd, args[0], &fields, jsonOutput)
		},
	}

	fields.bind(cmd)
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runUpdate(cmd *cobra.Command, id string, fields *calendarFields, jsonOutput bool) error {
	override := map[string]any{}
	fields.overlay(cmd, override)
	if len(override) == 0 {
		return fmt.Errorf("nothing to update; use --name, --color, --description or --echo-url")
	}

	api, err := loadAPI()
	if err != nil {
		return err
	}

	body, err := currentCalendar(api, id)
	if err != nil {
		return err
	}
	for key, value := range override {
		body[key] = value
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
	printCalendar(calendar)
	return nil
}

func currentCalendar(api *client.Client, id string) (map[string]any, error) {
	resp, err := api.Get("/calendars/" + id)
	if err != nil {
		return nil, err
	}
	var existing map[string]any
	if err := api.DecodeJSON(resp, &existing); err != nil {
		return nil, err
	}
	body := make(map[string]any, len(carriedFields))
	for _, key := range carriedFields {
		if value, ok := existing[key]; ok {
			body[key] = value
		}
	}
	return body, nil
}
