package caldav

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newSyncTokenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync-token <calendar-id>",
		Short: "Print the current sync token for a calendar",
		Long:  "Print the current CalDAV sync token (ctag) for a calendar. Useful for debugging sync issues.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSyncToken(args[0])
		},
	}
}

func runSyncToken(calendarID string) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	resp, err := api.Get("/calendars/" + calendarID)
	if err != nil {
		return err
	}

	var calendar map[string]any
	if err := api.DecodeJSON(resp, &calendar); err != nil {
		return err
	}

	syncToken, ok := calendar["sync_token"].(string)
	if !ok {
		ui.PrintWarn("Calendar has no sync token")
		return nil
	}

	ui.PrintSuccess("Sync token:")
	fmt.Println(syncToken)

	return nil
}
