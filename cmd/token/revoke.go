package token

import (
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newRevokeCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:     "revoke",
		Aliases: []string{"delete"},
		Short:   "Revoke personal API token",
		Long:    "Delete the active personal API token for the authenticated user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRevoke(jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runRevoke(jsonOutput bool) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	resp, err := api.Delete("/users/me/api-token")
	if err != nil {
		return err
	}

	var result map[string]any
	if err := api.DecodeJSON(resp, &result); err != nil {
		return err
	}

	if jsonOutput {
		return ui.PrintJSON(result)
	}

	ui.PrintSuccess("API token revoked")
	return nil
}
