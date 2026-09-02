package cmd

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

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

func init() {
	ui.Init()
}

var ErrShowHelp = fmt.Errorf("show help")
