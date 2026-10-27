package cmd

import (
	"github.com/FacileStudio/agenda-cli/cmd/token"
	"github.com/spf13/cobra"
)

// NewTokenCmd returns the token command group.
func NewTokenCmd() *cobra.Command {
	return token.NewTokenCmd()
}
