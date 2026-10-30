package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	var (
		force      bool
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a calendar",
		Long:  "Delete a calendar by ID. Requires write access.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(args[0], force, jsonOutput)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runDelete(id string, force, jsonOutput bool) error {
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

	if err := api.ExpectOK(resp); err != nil {
		return err
	}

	if jsonOutput {
		return ui.PrintJSON(map[string]bool{"deleted": true})
	}

	ui.PrintSuccess("Calendar deleted")
	return nil
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}
