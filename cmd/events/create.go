package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

type createEventOptions struct {
	calendarID  string
	title       string
	start       string
	end         string
	location    string
	description string
	allDay      bool
}

func newCreateCmd() *cobra.Command {
	var (
		opts       createEventOptions
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new event",
		Long:  "Create a new event in a calendar.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(opts, jsonOutput)
		},
	}

	cmd.Flags().StringVar(&opts.calendarID, "calendar", "", "Calendar ID")
	cmd.Flags().StringVar(&opts.title, "title", "", "Event title")
	cmd.Flags().StringVar(&opts.start, "start", "", "Start datetime (RFC3339)")
	cmd.Flags().StringVar(&opts.end, "end", "", "End datetime (RFC3339)")
	cmd.Flags().BoolVar(&opts.allDay, "all-day", false, "All-day event")
	cmd.Flags().StringVar(&opts.location, "location", "", "Event location")
	cmd.Flags().StringVar(&opts.description, "description", "", "Event description")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	cmd.MarkFlagRequired("calendar")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("start")
	cmd.MarkFlagRequired("end")

	return cmd
}

func runCreate(opts createEventOptions, jsonOutput bool) error {
	if err := requireFields(opts); err != nil {
		return err
	}

	startAt, err := parseDateTime(opts.start)
	if err != nil {
		return fmt.Errorf("invalid --start: %w", err)
	}

	endAt, err := parseDateTime(opts.end)
	if err != nil {
		return fmt.Errorf("invalid --end: %w", err)
	}

	opts.start = startAt
	opts.end = endAt

	api, err := loadAPI()
	if err != nil {
		return err
	}

	return createEvent(api, opts, jsonOutput)
}

func parseDateTime(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, layout := range formats {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.UTC().Format(time.RFC3339), nil
		}
	}
	return "", fmt.Errorf("expected RFC3339 or YYYY-MM-DD HH:MM")
}

func requireFields(opts createEventOptions) error {
	if err := requireField(opts.calendarID, "calendar"); err != nil {
		return err
	}
	if err := requireField(opts.title, "title"); err != nil {
		return err
	}
	if err := requireField(opts.start, "start"); err != nil {
		return err
	}
	return requireField(opts.end, "end")
}

func requireField(value, name string) error {
	if value == "" {
		return fmt.Errorf("--%s is required", name)
	}
	return nil
}

func createEvent(api *client.Client, opts createEventOptions, jsonOutput bool) error {
	body := map[string]any{
		"title":       opts.title,
		"start_at":    opts.start,
		"end_at":      opts.end,
		"is_all_day":  opts.allDay,
		"location":    opts.location,
		"description": opts.description,
	}

	resp, err := api.Post("/calendars/"+opts.calendarID+"/events", body)
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

	ui.PrintSuccess("Event created")
	ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
	ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))

	return nil
}
