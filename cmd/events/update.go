package events

import (
	"fmt"
	"strconv"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/spf13/cobra"
)

var carriedFields = []string{
	"title", "description", "location", "start_at", "end_at", "is_all_day",
	"recurrence_rule", "status", "conference_url", "conference_provider", "timezone",
}

func newUpdateCmd() *cobra.Command {
	var (
		occurrence string
		fields     eventFields
		calendarID string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an event",
		Long: "Update an existing event. Fields you do not pass keep their stored value; " +
			"pass an empty string to clear a text field. Moving between calendars is supported with --calendar. " +
			"--attendee and --reminder replace the stored list, --clear-attendees and --clear-reminders empty it. " +
			"With --occurrence YYYY-MM-DD only that occurrence of a recurring event changes " +
			"(--title, --description, --location, --start, --end, --status).",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if occurrence != "" {
				return runUpdateOccurrence(cmd, args[0], occurrence, &fields, jsonOutput)
			}
			return runUpdate(cmd, args[0], &fields, calendarID, jsonOutput)
		},
	}

	fields.bind(cmd)
	cmd.Flags().StringVar(&calendarID, "calendar", "", "Move event to another calendar")
	cmd.Flags().StringVar(&occurrence, "occurrence", "", "Change only this occurrence (YYYY-MM-DD or RFC3339 start)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runUpdate(cmd *cobra.Command, id string, fields *eventFields, calendarID string, jsonOutput bool) error {
	if cmd.Flags().NFlag() == 0 || (cmd.Flags().NFlag() == 1 && cmd.Flags().Changed("json")) {
		return fmt.Errorf("nothing to update: pass at least one field flag")
	}

	api, err := loadAPI()
	if err != nil {
		return err
	}
	existing, err := fetchEventMap(api, id)
	if err != nil {
		return err
	}

	body := carryFields(existing)
	fields.adoptStoredZone(cmd, existing)
	if err := fields.overlay(cmd, body); err != nil {
		return err
	}
	if calendarID != "" {
		calID, err := strconv.ParseInt(calendarID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid calendar id: %w", err)
		}
		body["calendar_id"] = calID
	}

	resp, err := api.Put("/events/"+id, body)
	if err != nil {
		return err
	}
	return printResult(api, resp, "Event updated", jsonOutput)
}

func fetchEventMap(api *client.Client, id string) (map[string]any, error) {
	resp, err := api.Get("/events/" + id)
	if err != nil {
		return nil, err
	}
	var existing map[string]any
	if err := api.DecodeJSON(resp, &existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func carryFields(existing map[string]any) map[string]any {
	body := make(map[string]any, len(carriedFields)+2)
	for _, key := range carriedFields {
		if value, ok := existing[key]; ok {
			body[key] = value
		}
	}
	return body
}
