package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newShareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "share <calendar-id> <email> <role>",
		Short: "Share a calendar with a user",
		Long:  "Share a calendar with a user by email. Role must be 'editor' or 'reader'.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			calendarID := args[0]
			email := args[1]
			role := args[2]

			if role != "editor" && role != "reader" {
				return fmt.Errorf("role must be 'editor' or 'reader'")
			}

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			api := client.NewClient(config)
			body := map[string]string{
				"email": email,
				"role":  role,
			}

			resp, err := api.Post("/calendars/"+calendarID+"/members", body)
			if err != nil {
				return err
			}

			var member map[string]interface{}
			if err := api.DecodeJSON(resp, &member); err != nil {
				return err
			}

			ui.PrintSuccess("Calendar shared")
			ui.PrintKeyValue("Email", email)
			ui.PrintKeyValue("Role", role)

			return nil
		},
	}
}
