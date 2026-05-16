package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
)

var trendingCmd = &cobra.Command{
	Use:          "trending",
	Short:        "Show currently trending anime",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt(keyPage)
		perPage, _ := cmd.Flags().GetInt("per-page")
		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")
		mediaType, _ := cmd.Flags().GetString("type")

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QueryTrending, map[string]any{
			"type":     strings.ToUpper(mediaType),
			keyPage:    page,
			keyPerPage: perPage,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return err
		}

		return printMedia(result.Media, asJSON, lang, noColor, mediaType)
	},
}

func init() {
	trendingCmd.Flags().Int(keyPage, 1, "page number")
	trendingCmd.Flags().Int("per-page", 20, "results per page (max 50)")
	rootCmd.AddCommand(trendingCmd)
}
