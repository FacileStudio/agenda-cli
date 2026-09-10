package caldav

import (
	"github.com/spf13/cobra"
)

// NewURLCmd returns the caldav url subcommand.
func NewURLCmd() *cobra.Command { return newURLCmd() }

// NewSyncTokenCmd returns the caldav sync-token subcommand.
func NewSyncTokenCmd() *cobra.Command { return newSyncTokenCmd() }

// NewTestCmd returns the caldav test subcommand.
func NewTestCmd() *cobra.Command { return newTestCmd() }
