package cmd

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

// NewRoot creates the root command for the agenda CLI.
func NewRoot() *cobra.Command {
	var version string

	root := &cobra.Command{
		Use:   "agenda",
		Short: "Agenda CLI — terminal client for Agenda",
		Long:  "A terminal client for the Facile Agenda calendar server.",
		Version: func() string {
			if version != "" {
				return version
			}
			return "dev"
		}(),
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	root.AddCommand(NewAuthCmd())
	root.AddCommand(NewCalendarsCmd())
	root.AddCommand(NewEventsCmd())
	root.AddCommand(NewCalDAVCmd())
	root.AddCommand(newConfigCmd())

	return root
}

// ErrShowHelp is returned when the user requests help.
var ErrShowHelp = fmt.Errorf("show help")

// init initializes the UI package.
func init() {
	ui.Init()
}
