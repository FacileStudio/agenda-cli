package events

import (
	"fmt"
	"time"

	"github.com/FacileStudio/agenda-cli/internal/caldav"
	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	var (
		ics bool
	)

	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show event details",
		Long:  "Show details for a specific event by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			api := client.NewClient(config)
			resp, err := api.Get("/events/" + id)
			if err != nil {
				return err
			}

			var event map[string]interface{}
			if err := api.DecodeJSON(resp, &event); err != nil {
				return err
			}

			if ics {
				calEvent := &caldav.Event{
					UID:         fmt.Sprintf("%v", event["id"]),
					Summary:     fmt.Sprintf("%v", event["title"]),
					Description: fmt.Sprintf("%v", event["description"]),
					Location:    fmt.Sprintf("%v", event["location"]),
					AllDay:      fmt.Sprintf("%v", event["all_day"]) == "true",
				}

				if startStr, ok := event["start"].(string); ok {
					if t, err := time.Parse(time.RFC3339, startStr); err == nil {
						calEvent.Start = t
					}
				}
				if endStr, ok := event["end"].(string); ok {
					if t, err := time.Parse(time.RFC3339, endStr); err == nil {
						calEvent.End = t
					}
				}

				icsData, err := caldav.RenderICS(calEvent)
				if err != nil {
					return err
				}
				fmt.Print(icsData)
				return nil
			}

			ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
			ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))
			ui.PrintKeyValue("Start", fmt.Sprintf("%v", event["start"]))
			ui.PrintKeyValue("End", fmt.Sprintf("%v", event["end"]))
			ui.PrintKeyValue("All Day", fmt.Sprintf("%v", event["all_day"]))
			ui.PrintKeyValue("Location", fmt.Sprintf("%v", event["location"]))
			ui.PrintKeyValue("Description", fmt.Sprintf("%v", event["description"]))
			ui.PrintKeyValue("Calendar", fmt.Sprintf("%v", event["calendar_id"]))

			return nil
		},
	}

	cmd.Flags().BoolVar(&ics, "ics", false, "Output as ICS format")

	return cmd
}
