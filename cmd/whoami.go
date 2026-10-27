package cmd

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Revoke the current session",
		Long:  "Call the server logout endpoint and remove the stored token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := client.Load()
			if err != nil {
				return err
			}

			if !config.IsAuthenticated() {
				ui.PrintWarn("Not logged in")
				return nil
			}

			api := client.NewClient(config)
			resp, err := api.Post("/auth/logout", nil)
			if err == nil {
				resp.Body.Close()
			}

			config.Token = ""
			config.UserID = ""
			config.Email = ""

			configPath := client.DefaultConfigPath()
			if err := config.Save(configPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			ui.PrintSuccess("Logged out")
			return nil
		},
	}
}

func newWhoamiCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "Show the current authenticated user",
		Long:  "Display the user_id and email from the stored config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWhoami(jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	return cmd
}

func runWhoami(jsonOutput bool) error {
	config, err := client.Load()
	if err != nil {
		return err
	}

	if !config.IsAuthenticated() {
		if jsonOutput {
			return fmt.Errorf("not logged in")
		}
		ui.PrintWarn("Not logged in")
		return nil
	}

	if config.UserID == "" || config.Email == "" {
		api := client.NewClient(config)
		if err := populateUser(api, config); err != nil {
			return err
		}
	}

	if jsonOutput {
		return ui.PrintJSON(map[string]string{
			"user_id":    config.UserID,
			"email":      config.Email,
			"server_url": config.ServerURL,
		})
	}

	ui.PrintKeyValue("User ID", config.UserID)
	ui.PrintKeyValue("Email", config.Email)
	ui.PrintKeyValue("Server", config.ServerURL)

	return nil
}
