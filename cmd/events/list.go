package events

import (
	"fmt"
	"time"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

type listFilter struct {
	calendarID string
	from       string
	to         string
	search     string
	timezone   string
	jsonOutput bool
}

func newListCmd() *cobra.Command {
	var filter listFilter

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List events",
		Long:  "List events in a calendar, optionally filtered by date range or search text.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(filter)
		},
	}

	cmd.Flags().StringVar(&filter.calendarID, "calendar", "", "Calendar ID")
	cmd.Flags().StringVar(&filter.from, "from", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&filter.to, "to", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&filter.search, "search", "", "Keep events whose title, description or location contain this text")
	cmd.Flags().StringVar(&filter.timezone, "timezone", "", "IANA zone used to expand recurring events")
	cmd.Flags().BoolVar(&filter.jsonOutput, "json", false, "Output in JSON format")

	cmd.MarkFlagRequired("calendar")

	return cmd
}

func runList(filter listFilter) error {
	if filter.calendarID == "" {
		return fmt.Errorf("--calendar is required")
	}

	api, err := loadAPI()
	if err != nil {
		return err
	}

	if _, err := loadZone(filter.timezone); err != nil {
		return err
	}
	params := buildQueryParams(filter.from, filter.to, filter.timezone)
	resp, err := api.GetWithQuery("/calendars/"+filter.calendarID+"/events", params)
	if err != nil {
		return err
	}

	var events []map[string]any
	if err := api.DecodeJSON(resp, &events); err != nil {
		return err
	}

	events = filterEvents(events, filter.search)

	if filter.jsonOutput {
		if events == nil {
			events = []map[string]any{}
		}
		return ui.PrintJSON(events)
	}

	printEvents(events)
	return nil
}

func buildQueryParams(from, to, timezone string) map[string]string {
	params := make(map[string]string)
	if from != "" {
		if _, err := time.Parse(time.RFC3339, from); err != nil {
			params["from"] = from + "T00:00:00Z"
		} else {
			params["from"] = from
		}
	}
	if to != "" {
		if _, err := time.Parse(time.RFC3339, to); err != nil {
			params["to"] = to + "T23:59:59Z"
		} else {
			params["to"] = to
		}
	}
	if timezone != "" {
		params["tz"] = timezone
	}
	return params
}

func filterEvents(events []map[string]any, search string) []map[string]any {
	if search == "" {
		return events
	}
	kept := make([]map[string]any, 0, len(events))
	for _, event := range events {
		if matchesSearch(event, search) {
			kept = append(kept, event)
		}
	}
	return kept
}
