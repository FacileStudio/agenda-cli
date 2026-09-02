package cmd

import (
	"fmt"
	"os"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

// NewAuthCmd returns the auth command group.
func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
		Long:  "Login, logout, and check authentication status.",
	}

	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newLogoutCmd())
	cmd.AddCommand(newWhoamiCmd())

	return cmd
}

// newLoginCmd returns the login subcommand.
func newLoginCmd() *cobra.Command {
	var password string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with the Agenda server",
		Long:  "Login using password authentication.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("read config: %w", err)
			}

			if config.ServerURL == "" {
				return fmt.Errorf("no server URL configured; run `agenda config server <url>` or set AGENDA_SERVER_URL")
			}

			ui.PrintInfo(fmt.Sprintf("Logging in to %s", config.ServerURL))

			if password == "" {
				return fmt.Errorf("password required; use --password flag")
			}

			if err := loginWithPassword(config, password); err != nil {
				return err
			}

			if err := config.Save(configPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			ui.PrintSuccess("Logged in successfully")
			return nil
		},
	}

	cmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	cmd.MarkFlagRequired("password")

	return cmd
}

// loginWithPassword authenticates the user with password credentials.
func loginWithPassword(config *client.Config, password string) error {
	client := client.NewClient(config)

	resp, err := client.Post("/auth/login", map[string]string{
		"email":    config.Email,
		"password": password,
	})
	if err != nil {
		return err
	}

	var authResp struct {
		UserID string `json:"user_id"`
		Token  string `json:"token"`
	}

	if err := client.DecodeJSON(resp, &authResp); err != nil {
		return err
	}

	config.UserID = authResp.UserID
	config.Token = authResp.Token

	return nil
}

// newLogoutCmd returns the logout subcommand.
func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Revoke the current session",
		Long:  "Call the server logout endpoint and remove the stored token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return fmt.Errorf("no config found at %s", configPath)
			}

			if !config.IsAuthenticated() {
				ui.PrintWarn("Not logged in")
				return nil
			}

			api := client.NewClient(config)
			resp, err := api.Post("/auth/logout", nil)
			if err != nil {
				return err
			}

			resp.Body.Close()

			config.Token = ""
			config.UserID = ""
			config.Email = ""

			if err := config.Save(configPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			ui.PrintSuccess("Logged out")
			return nil
		},
	}
}

// newWhoamiCmd returns the whoami subcommand.
func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show the current authenticated user",
		Long:  "Display the user_id and email from the stored config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				return fmt.Errorf("no config found at %s", configPath)
			}

			if !config.IsAuthenticated() {
				ui.PrintWarn("Not logged in")
				return nil
			}

			ui.PrintKeyValue("User ID", config.UserID)
			ui.PrintKeyValue("Email", config.Email)
			ui.PrintKeyValue("Server", config.ServerURL)

			return nil
		},
	}
}
