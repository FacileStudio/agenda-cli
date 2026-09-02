package events

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var (
		calendarID string
		from       string
		to         string
		search     string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List events",
		Long:  "List events in a calendar, optionally filtered by date range or search text.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if calendarID == "" {
				return fmt.Errorf("--calendar is required")
			}

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			api := client.NewClient(config)
			path := "/calendars/" + calendarID + "/events"

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

			req, err := http.NewRequest("GET", config.ServerURL+path, nil)
			if err != nil {
				return err
			}

			q := req.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			req.Header.Set("Accept", "application/json")
			if config.Token != "" {
				req.Header.Set("Authorization", config.AuthHeader())
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}

			var events []map[string]interface{}
			if err := api.DecodeJSON(resp, &events); err != nil {
				return err
			}

			ui.PrintSuccess(fmt.Sprintf("Found %d events", len(events)))
			for _, event := range events {
				ui.PrintKeyValue("ID", fmt.Sprintf("%v", event["id"]))
				ui.PrintKeyValue("Title", fmt.Sprintf("%v", event["title"]))
				ui.PrintKeyValue("Start", fmt.Sprintf("%v", event["start"]))
				ui.PrintKeyValue("End", fmt.Sprintf("%v", event["end"]))
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&calendarID, "calendar", "", "Calendar ID")
	cmd.Flags().StringVar(&from, "from", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&to, "to", "", "End date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&search, "search", "", "Search text in title/description")

	cmd.MarkFlagRequired("calendar")

	return cmd
}
