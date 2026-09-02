package caldav

import (
	"github.com/spf13/cobra"
)

func NewURLCmd() *cobra.Command { return newURLCmd() }
func NewSyncTokenCmd() *cobra.Command { return newSyncTokenCmd() }
func NewTestCmd() *cobra.Command { return newTestCmd() }
