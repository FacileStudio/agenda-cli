package events

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var (
		calendarID string
		title      string
		start      string
		end        string
		allDay     bool
		location   string
		description string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new event",
		Long:  "Create a new event in a calendar.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if calendarID == "" {
				return fmt.Errorf("--calendar is required")
			}
			if title == "" {
				return fmt.Errorf("--title is required")
			}
			if start == "" {
				return fmt.Errorf("--start is required")
			}
			if end == "" {
				return fmt.Errorf("--end is required")
			}

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			api := client.NewClient(config)
			body := map[string]interface{}{
				"calendar_id": calendarID,
				"title":       title,
				"start":       start,
				"end":         end,
				"all_day":     allDay,
				"location":    location,
				"description": description,
			}

			resp, err := api.Post("/calendars/"+calendarID+"/events", body)
			if err != nil {
				return err
			}

			var event map[string]interface{}
			if err := api.DecodeJSON(resp, &event); err != nil {
				return err
			}

			ui.PrintSuccess("Event created")
			ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
			ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))

			return nil
		},
	}

	cmd.Flags().StringVar(&calendarID, "calendar", "", "Calendar ID")
	cmd.Flags().StringVar(&title, "title", "", "Event title")
	cmd.Flags().StringVar(&start, "start", "", "Start datetime (RFC3339)")
	cmd.Flags().StringVar(&end, "end", "", "End datetime (RFC3339)")
	cmd.Flags().BoolVar(&allDay, "all-day", false, "All-day event")
	cmd.Flags().StringVar(&location, "location", "", "Event location")
	cmd.Flags().StringVar(&description, "description", "", "Event description")

	cmd.MarkFlagRequired("calendar")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("start")
	cmd.MarkFlagRequired("end")

	return cmd
}
