package calendars

import (
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command { return newListCmd() }
func NewCreateCmd() *cobra.Command { return newCreateCmd() }
func NewShowCmd() *cobra.Command { return newShowCmd() }
func NewUpdateCmd() *cobra.Command { return newUpdateCmd() }
func NewDeleteCmd() *cobra.Command { return newDeleteCmd() }
func NewShareCmd() *cobra.Command { return newShareCmd() }
func NewMembersCmd() *cobra.Command { return newMembersCmd() }
func NewUnshareCmd() *cobra.Command { return newUnshareCmd() }
