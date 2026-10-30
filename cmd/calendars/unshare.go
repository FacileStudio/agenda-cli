package calendars

import (
	"fmt"
	"strconv"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newUnshareCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "unshare <calendar-id> <member-id-or-email>",
		Short: "Remove a member from a calendar",
		Long:  "Remove a member's access from a calendar by user ID or email.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnshare(args[0], args[1], jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runUnshare(calendarID, memberTarget string, jsonOutput bool) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	memberID, err := resolveMemberID(api, calendarID, memberTarget)
	if err != nil {
		return err
	}

	resp, err := api.Delete("/calendars/" + calendarID + "/members/" + memberID)
	if err != nil {
		return err
	}
	if err := api.ExpectOK(resp); err != nil {
		return err
	}

	if jsonOutput {
		return ui.PrintJSON(map[string]bool{"deleted": true})
	}

	ui.PrintSuccess("Member removed")
	return nil
}

func resolveMemberID(api *client.Client, calendarID, memberTarget string) (string, error) {
	if _, err := strconv.ParseInt(memberTarget, 10, 64); err == nil {
		return memberTarget, nil
	}

	resp, err := api.Get("/calendars/" + calendarID + "/members")
	if err != nil {
		return "", err
	}

	var members []map[string]any
	if err := api.DecodeJSON(resp, &members); err != nil {
		return "", err
	}

	for _, m := range members {
		if fmt.Sprintf("%v", m["email"]) == memberTarget {
			return fmt.Sprintf("%v", m["user_id"]), nil
		}
	}

	return "", fmt.Errorf("member %q not found in calendar %s", memberTarget, calendarID)
}
