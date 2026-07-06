package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/tui"
)

// Running ani with no subcommand launches the TUI.
func init() {
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		g := getGlobalFlags(cmd)
		return tui.Start(newClient(), g.lang, g.mediaType, g.noColor)
	}
}
