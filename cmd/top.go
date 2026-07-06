package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
)

var topCmd = &cobra.Command{
	Use:          "top",
	Short:        "Show top anime by score",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		page, _, err := pageFlags(cmd)
		if err != nil {
			return err
		}
		g := getGlobalFlags(cmd)

		client := newClient()
		result, err := client.Query(context.Background(), anilist.QueryTop, map[string]any{
			keyType:    strings.ToUpper(g.mediaType),
			keyPage:    page,
			keyPerPage: clampPerPage(limit),
		})
		if err != nil {
			return err
		}

		return printMedia(cmd.OutOrStdout(), result.Media, g.asJSON, g.lang, g.noColor, g.mediaType)
	},
}

func init() {
	topCmd.Flags().Int("limit", 20, "number of results (max 50 per page)")
	topCmd.Flags().Int(keyPage, 1, "page number")
	rootCmd.AddCommand(topCmd)
}
