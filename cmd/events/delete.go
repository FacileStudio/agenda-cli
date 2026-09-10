package events

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an event",
		Long:  "Delete an event by ID.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(args[0], force)
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation")

	return cmd
}

func runDelete(id string, force bool) error {
	if !force && !confirm("Delete this event? [y/N] ") {
		return nil
	}

	api, err := loadAPI()
	if err != nil {
		return err
	}

	resp, err := api.Delete("/events/" + id)
	if err != nil {
		return err
	}

	resp.Body.Close()

	ui.PrintSuccess("Event deleted")
	return nil
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}
