package calendars

import (
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

			api, err := loadAPI()
			if err != nil {
				return err
			}

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
