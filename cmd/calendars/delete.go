package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a calendar",
		Long:  "Delete a calendar by ID. Requires write access.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(args[0], force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")

	return cmd
}

func runDelete(id string, force bool) error {
	if !force && !confirm("This will permanently delete the calendar. Continue? [y/N] ") {
		return nil
	}

	api, err := loadAPI()
	if err != nil {
		return err
	}

	resp, err := api.Delete("/calendars/" + id)
	if err != nil {
		return err
	}

	resp.Body.Close()

	ui.PrintSuccess("Calendar deleted")
	return nil
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}
