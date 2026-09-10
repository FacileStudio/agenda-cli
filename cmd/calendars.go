package cmd

import (
	"github.com/FacileStudio/agenda-cli/cmd/calendars"
	"github.com/spf13/cobra"
)

// NewCalendarsCmd returns the calendars command group.
func NewCalendarsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calendars",
		Short: "Manage calendars",
		Long:  "Create, list, update, and delete calendars.",
	}

	cmd.AddCommand(calendars.NewListCmd())
	cmd.AddCommand(calendars.NewCreateCmd())
	cmd.AddCommand(calendars.NewShowCmd())
	cmd.AddCommand(calendars.NewUpdateCmd())
	cmd.AddCommand(calendars.NewDeleteCmd())
	cmd.AddCommand(calendars.NewShareCmd())
	cmd.AddCommand(calendars.NewMembersCmd())
	cmd.AddCommand(calendars.NewUnshareCmd())

	return cmd
}
