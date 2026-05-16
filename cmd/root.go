package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	keyPage    = "page"
	keyPerPage = "perPage"
	keyType    = "type"
)

var rootCmd = &cobra.Command{
	Use:          "ani",
	Short:        "browse and search AniList anime from your terminal",
	SilenceUsage: true,
}

// Execute is the entry point called from main.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().Bool("json", false, "output as JSON")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable color output")
	rootCmd.PersistentFlags().String("lang", "english", "title language: romaji, english, or native")
	rootCmd.PersistentFlags().String("type", "anime", "media type: anime or manga")

	// Respect NO_COLOR env variable.
	if os.Getenv("NO_COLOR") != "" {
		_ = rootCmd.PersistentFlags().Set("no-color", "true")
	}
}
