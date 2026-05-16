package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/pavelsimo/ani/internal/tui"
)

func init() {
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")
		mediaType, _ := cmd.Flags().GetString("type")
		return tui.Start(anilist.New(), lang, mediaType)
	}
}
