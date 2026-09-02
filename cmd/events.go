package cmd

import (
	"github.com/FacileStudio/agenda-cli/cmd/events"
	"github.com/spf13/cobra"
)

func NewEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Manage events",
		Long:  "Create, list, show, update, and delete events.",
	}

	cmd.AddCommand(events.NewListCmd())
	cmd.AddCommand(events.NewShowCmd())
	cmd.AddCommand(events.NewCreateCmd())
	cmd.AddCommand(events.NewUpdateCmd())
	cmd.AddCommand(events.NewDeleteCmd())

	return cmd
}
