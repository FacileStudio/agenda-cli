package events

import (
	"fmt"
	"strconv"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

type updateEventOptions struct {
	title       string
	start       string
	end         string
	allDay      bool
	location    string
	description string
	calendarID  string
}

func newUpdateCmd() *cobra.Command {
	var (
		opts       updateEventOptions
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an event",
		Long:  "Update an existing event. Moving between calendars is supported with --calendar.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(args[0], opts, jsonOutput)
		},
	}

	cmd.Flags().StringVar(&opts.title, "title", "", "New event title")
	cmd.Flags().StringVar(&opts.start, "start", "", "New start datetime (RFC3339)")
	cmd.Flags().StringVar(&opts.end, "end", "", "New end datetime (RFC3339)")
	cmd.Flags().BoolVar(&opts.allDay, "all-day", false, "All-day event")
	cmd.Flags().StringVar(&opts.location, "location", "", "New event location")
	cmd.Flags().StringVar(&opts.description, "description", "", "New event description")
	cmd.Flags().StringVar(&opts.calendarID, "calendar", "", "Move event to another calendar")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runUpdate(id string, opts updateEventOptions, jsonOutput bool) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	existing, err := fetchEventMap(api, id)
	if err != nil {
		return err
	}

	body, err := buildEventPatch(existing, opts)
	if err != nil {
		return err
	}

	resp, err := api.Put("/events/"+id, body)
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

	ui.PrintSuccess("Event updated")
	ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
	ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))

	return nil
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

func buildEventPatch(existing map[string]any, opts updateEventOptions) (map[string]any, error) {
	body := map[string]any{
		"title":       existing["title"],
		"start_at":    existing["start_at"],
		"end_at":      existing["end_at"],
		"is_all_day":  existing["is_all_day"],
		"location":    existing["location"],
		"description": existing["description"],
	}

	if err := applyDateOverrides(body, opts); err != nil {
		return nil, err
	}
	applyTextOverrides(body, opts)

	if opts.calendarID != "" {
		calID, err := strconv.ParseInt(opts.calendarID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid calendar id: %w", err)
		}
		body["calendar_id"] = calID
	}
	return body, nil
}

func applyDateOverrides(body map[string]any, opts updateEventOptions) error {
	if opts.start != "" {
		startAt, err := parseDateTime(opts.start)
		if err != nil {
			return fmt.Errorf("invalid --start: %w", err)
		}
		body["start_at"] = startAt
	}
	if opts.end != "" {
		endAt, err := parseDateTime(opts.end)
		if err != nil {
			return fmt.Errorf("invalid --end: %w", err)
		}
		body["end_at"] = endAt
	}
	if opts.allDay {
		body["is_all_day"] = true
	}
	return nil
}

func applyTextOverrides(body map[string]any, opts updateEventOptions) {
	if opts.title != "" {
		body["title"] = opts.title
	}
	if opts.location != "" {
		body["location"] = opts.location
	}
	if opts.description != "" {
		body["description"] = opts.description
	}
}
