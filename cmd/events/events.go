package events

import (
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command { return newListCmd() }
func NewShowCmd() *cobra.Command { return newShowCmd() }
func NewCreateCmd() *cobra.Command { return newCreateCmd() }
func NewUpdateCmd() *cobra.Command { return newUpdateCmd() }
func NewDeleteCmd() *cobra.Command { return newDeleteCmd() }
