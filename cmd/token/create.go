package token

import (
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var (
		name       string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create personal API token",
		Long:  "Create or replace the personal API token for the authenticated user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(name, jsonOutput)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Token label or name")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runCreate(name string, jsonOutput bool) error {
	api, err := loadAPI()
	if err != nil {
		return err
	}

	body := map[string]string{}
	if name != "" {
		body["name"] = name
	}

	resp, err := api.Post("/users/me/api-token", body)
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

	printCreated(result)
	return nil
}

func printCreated(result map[string]any) {
	ui.PrintSuccess("API token created")
	if token, ok := result["token"].(string); ok {
		ui.PrintKeyValue("Token", token)
	}
	if n, ok := result["name"].(string); ok && n != "" {
		ui.PrintKeyValue("Name", n)
	}
	if createdAt, ok := result["created_at"].(string); ok && createdAt != "" {
		ui.PrintKeyValue("Created At", createdAt)
	}
}
