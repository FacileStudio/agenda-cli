package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var (
		fields     calendarFields
		spaceID    int64
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new calendar",
		Long:  "Create a new calendar. The name is the argument or --name.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				cmd.Flags().Set("name", args[0])
			}
			return runCreate(cmd, &fields, spaceID, jsonOutput)
		},
	}

	fields.bind(cmd)
	cmd.Flags().Int64Var(&spaceID, "space", 0, "Scope the calendar to this space id")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runCreate(cmd *cobra.Command, fields *calendarFields, spaceID int64, jsonOutput bool) error {
	body := map[string]any{}
	fields.overlay(cmd, body)
	if _, ok := body["name"]; !ok {
		return fmt.Errorf("a calendar needs a name: pass it as the argument or with --name")
	}
	if spaceID != 0 {
		body["space_id"] = spaceID
	}

	api, err := loadAPI()
	if err != nil {
		return err
	}

	resp, err := api.Post("/calendars", body)
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

	ui.PrintSuccess("Calendar created")
	printCalendar(calendar)
	return nil
}

func printCalendar(calendar map[string]any) {
	for _, row := range [][2]string{
		{"ID", "id"}, {"Name", "name"}, {"Color", "color"}, {"Description", "description"},
		{"Echo URL", "echo_url"}, {"Role", "role"}, {"Owner", "owner_id"},
	} {
		if v, ok := calendar[row[1]]; ok && v != nil {
			ui.PrintKeyValue(row[0], fmt.Sprintf("%v", v))
		}
	}
}
