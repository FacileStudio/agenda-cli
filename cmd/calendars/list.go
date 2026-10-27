package calendars

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all calendars",
		Long:  "List all calendars visible to the authenticated user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runList(jsonOutput bool) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}
	resp, err := api.Get("/calendars")
	if err != nil {
		return err
	}

	var calendars []map[string]any
	if err := api.DecodeJSON(resp, &calendars); err != nil {
		return err
	}

	if jsonOutput {
		if calendars == nil {
			calendars = []map[string]any{}
		}
		return ui.PrintJSON(calendars)
	}

	ui.PrintSuccess(fmt.Sprintf("Found %d calendars", len(calendars)))
	for _, cal := range calendars {
		ui.PrintKeyValue("ID", fmt.Sprintf("%v", cal["id"]))
		ui.PrintKeyValue("Name", fmt.Sprintf("%v", cal["name"]))
		ui.PrintKeyValue("Color", fmt.Sprintf("%v", cal["color"]))
		fmt.Println()
	}

	return nil
}
