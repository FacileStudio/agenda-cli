package cmd

import (
	"github.com/FacileStudio/agenda-cli/cmd/caldav"
	"github.com/spf13/cobra"
)

// NewCalDAVCmd returns the caldav command group.
func NewCalDAVCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "caldav",
		Short: "CalDAV utilities",
		Long:  "CalDAV setup, debugging, and connectivity helpers.",
	}

	cmd.AddCommand(caldav.NewURLCmd())
	cmd.AddCommand(caldav.NewSyncTokenCmd())
	cmd.AddCommand(caldav.NewTestCmd())

	return cmd
}
