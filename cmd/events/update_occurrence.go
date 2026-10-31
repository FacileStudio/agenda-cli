package events

import (
	"fmt"
	"slices"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/spf13/cobra"
)

var seriesOnlyFlags = []string{
	"recurrence", "all-day", "conference-url", "conference-provider", "attendee", "clear-attendees",
	"reminder", "clear-reminders", "color", "calendar",
}

func runUpdateOccurrence(cmd *cobra.Command, id, when string, fields *eventFields, jsonOutput bool) error {
	if err := checkOccurrenceFlags(cmd, fields); err != nil {
		return err
	}
	api, err := loadAPI()
	if err != nil {
		return err
	}
	body, err := occurrenceBody(cmd, api, id, fields)
	if err != nil {
		return err
	}
	resp, err := api.Put(occurrencePath(id, when, ""), body)
	if err != nil {
		return err
	}
	return printResult(api, resp, "Occurrence updated", jsonOutput)
}

func checkOccurrenceFlags(cmd *cobra.Command, fields *eventFields) error {
	for _, name := range seriesOnlyFlags {
		if cmd.Flags().Changed(name) {
			return fmt.Errorf("--%s applies to the whole series: drop it or drop --occurrence", name)
		}
	}
	if fields.status != "" && !slices.Contains([]string{"confirmed", "tentative"}, fields.status) {
		return fmt.Errorf("an occurrence --status is confirmed or tentative: use cancel-occurrence to cancel it")
	}
	return nil
}

func occurrenceBody(cmd *cobra.Command, api *client.Client, id string, fields *eventFields) (map[string]any, error) {
	if !cmd.Flags().Changed("timezone") {
		existing, err := fetchEventMap(api, id)
		if err != nil {
			return nil, err
		}
		fields.adoptStoredZone(cmd, existing)
	}
	body := map[string]any{}
	if err := fields.overlay(cmd, body); err != nil {
		return nil, err
	}
	delete(body, "timezone")
	if len(body) == 0 {
		return nil, fmt.Errorf("nothing to update: pass --title, --description, --location, --start, --end or --status")
	}
	return body, nil
}
