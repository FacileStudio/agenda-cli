package events

import (
	"fmt"
	"time"

	"github.com/FacileStudio/agenda-cli/internal/caldav"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	var ics bool

	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show event details",
		Long:  "Show details for a specific event by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(args[0], ics)
		},
	}

	cmd.Flags().BoolVar(&ics, "ics", false, "Output as ICS format")

	return cmd
}

func runShow(id string, ics bool) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	resp, err := api.Get("/events/" + id)
	if err != nil {
		return err
	}

	var event map[string]any
	if err := api.DecodeJSON(resp, &event); err != nil {
		return err
	}

	if ics {
		data, err := renderICS(event)
		if err != nil {
			return err
		}
		fmt.Print(data)
		return nil
	}

	printEvent(event)
	return nil
}

func renderICS(event map[string]any) (string, error) {
	calEvent := &caldav.Event{
		UID:         fmt.Sprintf("%v", event["id"]),
		Summary:     fmt.Sprintf("%v", event["title"]),
		Description: fmt.Sprintf("%v", event["description"]),
		Location:    fmt.Sprintf("%v", event["location"]),
		AllDay:      fmt.Sprintf("%v", event["all_day"]) == "true",
		Start:       parseTime(event["start"]),
		End:         parseTime(event["end"]),
	}
	return caldav.RenderICS(calEvent)
}

func parseTime(v any) time.Time {
	if s, ok := v.(string); ok {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func printEvent(event map[string]any) {
	ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
	ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))
	ui.PrintKeyValue("Start", fmt.Sprintf("%v", event["start"]))
	ui.PrintKeyValue("End", fmt.Sprintf("%v", event["end"]))
	ui.PrintKeyValue("All Day", fmt.Sprintf("%v", event["all_day"]))
	ui.PrintKeyValue("Location", fmt.Sprintf("%v", event["location"]))
	ui.PrintKeyValue("Description", fmt.Sprintf("%v", event["description"]))
	ui.PrintKeyValue("Calendar", fmt.Sprintf("%v", event["calendar_id"]))
}
