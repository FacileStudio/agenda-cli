package cmd

import (
	"github.com/FacileStudio/agenda-cli/cmd/events"
	"github.com/spf13/cobra"
)

// NewEventsCmd returns the events command group.
func NewEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Manage events",
		Long:  "Create, list, show, update, and delete events, and cancel or change single occurrences of a recurring event.",
	}

	cmd.AddCommand(events.NewListCmd())
	cmd.AddCommand(events.NewShowCmd())
	cmd.AddCommand(events.NewCreateCmd())
	cmd.AddCommand(events.NewUpdateCmd())
	cmd.AddCommand(events.NewDeleteCmd())
	cmd.AddCommand(events.NewCancelOccurrenceCmd())
	cmd.AddCommand(events.NewRestoreOccurrenceCmd())

	return cmd
}
