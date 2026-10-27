package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newMembersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "members <calendar-id>",
		Short: "List calendar members",
		Long:  "List all members of a calendar.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			calendarID := args[0]

			api, err := loadAPI()
			if err != nil {
				return err
			}

			resp, err := api.Get("/calendars/" + calendarID + "/members")
			if err != nil {
				return err
			}

			var members []map[string]any
			if err := api.DecodeJSON(resp, &members); err != nil {
				return err
			}

			ui.PrintSuccess(fmt.Sprintf("Found %d members", len(members)))
			for _, member := range members {
				ui.PrintKeyValue("Email", fmt.Sprintf("%v", member["email"]))
				ui.PrintKeyValue("Role", fmt.Sprintf("%v", member["role"]))
				fmt.Println()
			}

			return nil
		},
	}
}
