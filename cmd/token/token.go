package token

import (
	"github.com/spf13/cobra"
)

// NewTokenCmd returns the token command group.
func NewTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Manage personal API tokens",
		Long:  "Check status, create, or revoke personal API tokens.",
	}

	cmd.AddCommand(newStatusCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newRevokeCmd())

	return cmd
}
