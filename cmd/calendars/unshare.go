package calendars

import (
	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newUnshareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unshare <calendar-id> <member-email>",
		Short: "Remove a member from a calendar",
		Long:  "Remove a member's access from a calendar.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			calendarID := args[0]
			memberEmail := args[1]

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return err
			}

			api := client.NewClient(config)
			resp, err := api.Delete("/calendars/" + calendarID + "/members/" + memberEmail)
			if err != nil {
				return err
			}

			resp.Body.Close()

			ui.PrintSuccess("Member removed")
			return nil
		},
	}
}
