package calendars

import (
	"github.com/spf13/cobra"
)

// NewListCmd returns the calendars list subcommand.
func NewListCmd() *cobra.Command { return newListCmd() }

// NewCreateCmd returns the calendars create subcommand.
func NewCreateCmd() *cobra.Command { return newCreateCmd() }

// NewShowCmd returns the calendars show subcommand.
func NewShowCmd() *cobra.Command { return newShowCmd() }

// NewUpdateCmd returns the calendars update subcommand.
func NewUpdateCmd() *cobra.Command { return newUpdateCmd() }

// NewDeleteCmd returns the calendars delete subcommand.
func NewDeleteCmd() *cobra.Command { return newDeleteCmd() }

// NewShareCmd returns the calendars share subcommand.
func NewShareCmd() *cobra.Command { return newShareCmd() }

// NewMembersCmd returns the calendars members subcommand.
func NewMembersCmd() *cobra.Command { return newMembersCmd() }

// NewUnshareCmd returns the calendars unshare subcommand.
func NewUnshareCmd() *cobra.Command { return newUnshareCmd() }
