package cmd

import (
	"fmt"
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
		Long:  "Login, logout, and check authentication status.",
	}

	cmd.AddCommand(newLoginCmd())
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
		Long:  "Login using browser-based OIDC flow.",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := client.DefaultConfigPath()
			config, err := client.Load()
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			if config.ServerURL == "" {
				return fmt.Errorf("no server URL configured; run `agenda config server <url>` or set AGENDA_SERVER_URL")
			}

			ui.PrintInfo(fmt.Sprintf("Logging in to %s", config.ServerURL))

			if err := loginWithOIDC(config, noBrowser); err != nil {
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

func loginWithOIDC(config *client.Config, noBrowser bool) error {
	api := client.NewClient(config)

	authConfig, err := api.AuthConfig()
	if err != nil {
		return fmt.Errorf("failed to get auth config: %w", err)
	}
	if !authConfig.OIDCEnabled {
		return fmt.Errorf("server does not support OIDC login")
	}

	authCode, err := obtainAuthCode(config.ServerURL, noBrowser)
	if err != nil {
		return err
	}

	token, err := exchangeToken(api, authCode)
	if err != nil {
		return err
	}

	config.Token = *token
	return populateUser(api, config)
}

func obtainAuthCode(serverURL string, noBrowser bool) (string, error) {
	listener, port, err := loopback.Listen()
	if err != nil {
		return "", fmt.Errorf("failed to start listener: %w", err)
	}
	defer listener.Close()

	state := loopback.RandomState()
	authURL := fmt.Sprintf("%s/auth/oidc?flow=cli&port=%d&cli_state=%s", serverURL, port, state)

	ui.PrintInfo("Opening browser for authentication...")
	if !noBrowser && !loopback.OpenBrowser(authURL) {
		ui.PrintWarn("Could not open browser automatically")
	}
	ui.PrintInfo(fmt.Sprintf("If browser did not open, visit: %s", authURL))
	ui.PrintInfo("Waiting for authentication...")

	authCode, err := loopback.WaitForCode(listener, state, 5*time.Minute)
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}
	return authCode, nil
}

func exchangeToken(api *client.Client, authCode string) (*string, error) {
	token, err := api.Exchange(authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}
	if token == nil || *token == "" {
		return nil, fmt.Errorf("received empty token from server")
	}
	return token, nil
}

func populateUser(api *client.Client, config *client.Config) error {
	resp, err := api.Get("/users/me")
	if err != nil {
		resp, err = api.Get("/me")
	}
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var env struct {
		User struct {
			ID    any    `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		ID    any    `json:"id"`
		Email string `json:"email"`
	}
	if err := api.DecodeJSON(resp, &env); err != nil {
		return fmt.Errorf("failed to decode user info: %w", err)
	}
	applyUserEnvelope(config, env.User.Email, env.User.ID, env.Email, env.ID)
	return nil
}

func applyUserEnvelope(config *client.Config, userEmail string, userID any, rootEmail string, rootID any) {
	if userEmail != "" {
		config.Email = userEmail
		config.UserID = fmt.Sprintf("%v", userID)
		return
	}
	config.Email = rootEmail
	if rootID != nil {
		config.UserID = fmt.Sprintf("%v", rootID)
	}
}
