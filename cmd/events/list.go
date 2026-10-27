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
	cmd.Flags().StringVar(&filter.search, "search", "", "Search text in title/description")
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

	params := buildQueryParams(filter.from, filter.to, filter.search)
	resp, err := api.GetWithQuery("/calendars/"+filter.calendarID+"/events", params)
	if err != nil {
		return err
	}

	var events []map[string]any
	if err := api.DecodeJSON(resp, &events); err != nil {
		return err
	}

	if filter.jsonOutput {
		if events == nil {
			events = []map[string]any{}
		}
		return ui.PrintJSON(events)
	}

	printEvents(events)
	return nil
}

func printEvents(events []map[string]any) {
	ui.PrintSuccess(fmt.Sprintf("Found %d events", len(events)))
	for _, event := range events {
		ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
		ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))
		ui.PrintKeyValue("Start", fmt.Sprintf("%v", event["start"]))
		ui.PrintKeyValue("End", fmt.Sprintf("%v", event["end"]))
		fmt.Println()
	}
}

func buildQueryParams(from, to, search string) map[string]string {
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
	if search != "" {
		params["search"] = search
	}
	return params
}
