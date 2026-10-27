package token

import (
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check personal API token status",
		Long:  "Inspect whether the authenticated user has an active API token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runStatus(jsonOutput bool) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	resp, err := api.Get("/users/me/api-token")
	if err != nil {
		return err
	}

	var status map[string]any
	if err := api.DecodeJSON(resp, &status); err != nil {
		return err
	}

	if jsonOutput {
		return ui.PrintJSON(status)
	}

	printStatus(status)
	return nil
}

func printStatus(status map[string]any) {
	hasToken, _ := status["has_token"].(bool)
	if !hasToken {
		ui.PrintInfo("No active API token")
		return
	}

	ui.PrintSuccess("Active API token")
	if name, ok := status["name"].(string); ok && name != "" {
		ui.PrintKeyValue("Name", name)
	}
	if createdAt, ok := status["created_at"].(string); ok && createdAt != "" {
		ui.PrintKeyValue("Created At", createdAt)
	}
}
