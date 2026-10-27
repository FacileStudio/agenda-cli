package cmd

import (
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

var version string

// NewRoot creates the root command for the agenda CLI.
func NewRoot() *cobra.Command {
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
	root.SetVersionTemplate("{{.Name}} {{.Version}}\n")

	root.AddCommand(newLoginCmd())
	root.AddCommand(newLogoutCmd())
	root.AddCommand(newWhoamiCmd())
	root.AddCommand(NewAuthCmd())
	root.AddCommand(NewCalendarsCmd())
	root.AddCommand(NewEventsCmd())
	root.AddCommand(NewCalDAVCmd())
	root.AddCommand(newConfigCmd())

	ui.Init()

	return root
}
