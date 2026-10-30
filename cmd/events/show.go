package events

import (
	"fmt"
	"time"

	"github.com/FacileStudio/agenda-cli/internal/caldav"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	var (
		ics        bool
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show event details",
		Long:  "Show details for a specific event by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(args[0], ics, jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&ics, "ics", false, "Output as ICS format")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runShow(id string, ics, jsonOutput bool) error {
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

	if jsonOutput {
		return ui.PrintJSON(event)
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
		UID:         firstText(event["uid"], event["id"]),
		Summary:     text(event["title"]),
		Description: text(event["description"]),
		Location:    text(event["location"]),
		AllDay:      text(event["is_all_day"]) == "true",
		Start:       parseTime(event["start_at"]),
		End:         parseTime(event["end_at"]),
		Recurrence:  text(event["recurrence_rule"]),
		Status:      text(event["status"]),
		URL:         text(event["conference_url"]),
		CreatedAt:   parseTime(event["created_at"]),
		UpdatedAt:   parseTime(event["updated_at"]),
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

func firstText(values ...any) string {
	for _, v := range values {
		if s := text(v); s != "" {
			return s
		}
	}
	return ""
}
