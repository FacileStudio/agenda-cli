package events

import (
	"fmt"

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
	var opts updateEventOptions

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an event",
		Long:  "Update an existing event. Moving between calendars is supported with --calendar.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(args[0], opts)
		},
	}

	cmd.Flags().StringVar(&opts.title, "title", "", "New event title")
	cmd.Flags().StringVar(&opts.start, "start", "", "New start datetime (RFC3339)")
	cmd.Flags().StringVar(&opts.end, "end", "", "New end datetime (RFC3339)")
	cmd.Flags().BoolVar(&opts.allDay, "all-day", false, "All-day event")
	cmd.Flags().StringVar(&opts.location, "location", "", "New event location")
	cmd.Flags().StringVar(&opts.description, "description", "", "New event description")
	cmd.Flags().StringVar(&opts.calendarID, "calendar", "", "Move event to another calendar")

	return cmd
}

func runUpdate(id string, opts updateEventOptions) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	body := buildEventPatch(opts)
	if len(body) == 0 {
		return fmt.Errorf("nothing to update")
	}

	resp, err := api.Put("/events/"+id, body)
	if err != nil {
		return err
	}

	var event map[string]any
	if err := api.DecodeJSON(resp, &event); err != nil {
		return err
	}

	ui.PrintSuccess("Event updated")
	ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
	ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))

	return nil
}

func buildEventPatch(opts updateEventOptions) map[string]any {
	body := make(map[string]any)
	if opts.title != "" {
		body["title"] = opts.title
	}
	if opts.start != "" {
		body["start"] = opts.start
	}
	if opts.end != "" {
		body["end"] = opts.end
	}
	body["all_day"] = opts.allDay
	if opts.location != "" {
		body["location"] = opts.location
	}
	if opts.description != "" {
		body["description"] = opts.description
	}
	if opts.calendarID != "" {
		body["calendar_id"] = opts.calendarID
	}
	return body
}
