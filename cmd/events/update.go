package events

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var (
		title       string
		start       string
		end         string
		allDay      bool
		location    string
		description string
		calendarID  string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an event",
		Long:  "Update an existing event. Moving between calendars is supported with --calendar.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			body := make(map[string]interface{})
			if title != "" {
				body["title"] = title
			}
			if start != "" {
				body["start"] = start
			}
			if end != "" {
				body["end"] = end
			}
			body["all_day"] = allDay
			if location != "" {
				body["location"] = location
			}
			if description != "" {
				body["description"] = description
			}
			if calendarID != "" {
				body["calendar_id"] = calendarID
			}

			if len(body) == 0 {
				return fmt.Errorf("nothing to update")
			}

			api := client.NewClient(config)
			resp, err := api.Put("/events/"+id, body)
			if err != nil {
				return err
			}

			var event map[string]interface{}
			if err := api.DecodeJSON(resp, &event); err != nil {
				return err
			}

			ui.PrintSuccess("Event updated")
			ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
			ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))

			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New event title")
	cmd.Flags().StringVar(&start, "start", "", "New start datetime (RFC3339)")
	cmd.Flags().StringVar(&end, "end", "", "New end datetime (RFC3339)")
	cmd.Flags().BoolVar(&allDay, "all-day", false, "All-day event")
	cmd.Flags().StringVar(&location, "location", "", "New event location")
	cmd.Flags().StringVar(&description, "description", "", "New event description")
	cmd.Flags().StringVar(&calendarID, "calendar", "", "Move event to another calendar")

	return cmd
}
