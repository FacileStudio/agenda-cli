package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newShareCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "share <calendar-id> <email> <role>",
		Short: "Share a calendar with a user",
		Long:  "Share a calendar with a user by email. Role must be 'editor' or 'reader'.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShare(args[0], args[1], args[2], jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runShare(calendarID, email, role string, jsonOutput bool) error {
	if role != "editor" && role != "reader" {
		return fmt.Errorf("role must be 'editor' or 'reader'")
	}

	api, err := loadAPI()
	if err != nil {
		return err
	}

	body := map[string]string{
		"email": email,
		"role":  role,
	}

	resp, err := api.Post("/calendars/"+calendarID+"/members", body)
	if err != nil {
		return err
	}

	var member map[string]any
	if err := api.DecodeJSON(resp, &member); err != nil {
		return err
	}

	if jsonOutput {
		return ui.PrintJSON(member)
	}

	ui.PrintSuccess("Calendar shared")
	ui.PrintKeyValue("Email", email)
	ui.PrintKeyValue("Role", role)

	return nil
}
