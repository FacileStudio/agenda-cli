package cmd

import (
	"fmt"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Get or set configuration values.",
	}

	cmd.AddCommand(newConfigServerCmd())
	cmd.AddCommand(newConfigPathCmd())

	return cmd
}

func newConfigServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server <url>",
		Short: "Set the Agenda server URL",
		Long:  "Set the Agenda server URL in the config file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]

			configPath := client.DefaultConfigPath()
			config := &client.Config{}
			if err := config.Read(configPath); err != nil {
				config = &client.Config{}
			}

			config.ServerURL = url

			if err := config.Save(configPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			ui.PrintSuccess("Server URL saved")
			ui.PrintDim(configPath)

			return nil
		},
	}
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		Long:  "Print the path to the config file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := client.DefaultConfigPath()
			fmt.Println(path)
			return nil
		},
	}
}
