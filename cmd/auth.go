package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/loopback"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

// NewAuthCmd returns the auth command group.
func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication commands",
		Long:  "Logout, and check authentication status.",
	}

	cmd.AddCommand(newLogoutCmd())
	cmd.AddCommand(newWhoamiCmd())

	return cmd
}

// newLoginCmd returns the login subcommand.
func newLoginCmd() *cobra.Command {
	var noBrowser bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with the Agenda server",
		Long:  "Login using device code flow.",
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

			if err := loginWithDeviceCode(config, noBrowser); err != nil {
				return err
			}

			if err := config.Save(configPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			ui.PrintSuccess("Logged in successfully")
			return nil
		},
	}

	cmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Do not open browser automatically")

	return cmd
}

// loginWithDeviceCode authenticates the user using the device code flow.
func loginWithDeviceCode(config *client.Config, noBrowser bool) error {
	api := client.NewClient(config)

	authConfig, err := api.AuthConfig()
	if err != nil {
		return fmt.Errorf("failed to get auth config: %w", err)
	}
	if !authConfig.OIDCEnabled {
		return fmt.Errorf("server does not support device code login")
	}

	deviceCode, err := obtainDeviceCode(config.ServerURL, noBrowser)
	if err != nil {
		return err
	}

	token, err := exchangeToken(api, deviceCode)
	if err != nil {
		return err
	}

	config.Token = *token
	return populateUser(api, config)
}

func obtainDeviceCode(serverURL string, noBrowser bool) (string, error) {
	listener, port, err := loopback.Listen()
	if err != nil {
		return "", fmt.Errorf("failed to start listener: %w", err)
	}
	defer listener.Close()

	state := loopback.RandomState()
	authURL := fmt.Sprintf("%s/auth/device?flow=cli&port=%d&cli_state=%s", serverURL, port, state)

	ui.PrintInfo("Opening browser for authentication...")
	if !noBrowser && !loopback.OpenBrowser(authURL) {
		ui.PrintWarn("Could not open browser automatically")
	}
	ui.PrintInfo(fmt.Sprintf("If browser did not open, visit: %s", authURL))
	ui.PrintInfo("Waiting for authentication...")

	deviceCode, err := loopback.WaitForCode(listener, state, 5*time.Minute)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}
	return deviceCode, nil
}

func exchangeToken(api *client.Client, deviceCode string) (*string, error) {
	token, err := api.Exchange(deviceCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange device code: %w", err)
	}
	if token == nil || *token == "" {
		return nil, fmt.Errorf("received empty token from server")
	}
	return token, nil
}

func populateUser(api *client.Client, config *client.Config) error {
	resp, err := api.Get("/me")
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var meResp struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
	}

	if err := api.DecodeJSON(resp, &meResp); err != nil {
		return fmt.Errorf("failed to decode user info: %w", err)
	}

	config.UserID = meResp.UserID
	config.Email = meResp.Email

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
