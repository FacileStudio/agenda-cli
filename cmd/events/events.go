package events

import (
	"github.com/spf13/cobra"
)

// NewListCmd returns the events list subcommand.
func NewListCmd() *cobra.Command { return newListCmd() }

// NewShowCmd returns the events show subcommand.
func NewShowCmd() *cobra.Command { return newShowCmd() }

// NewCreateCmd returns the events create subcommand.
func NewCreateCmd() *cobra.Command { return newCreateCmd() }

// NewUpdateCmd returns the events update subcommand.
func NewUpdateCmd() *cobra.Command { return newUpdateCmd() }

// NewDeleteCmd returns the events delete subcommand.
func NewDeleteCmd() *cobra.Command { return newDeleteCmd() }
