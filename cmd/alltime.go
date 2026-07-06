package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
)

var alltimeCmd = &cobra.Command{
	Use:          "all-time",
	Short:        "Show most popular anime of all time",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		page, perPage, err := pageFlags(cmd)
		if err != nil {
			return err
		}
		g := getGlobalFlags(cmd)

		client := newClient()
		result, err := client.Query(context.Background(), anilist.QueryAllTime, map[string]any{
			keyType:    strings.ToUpper(g.mediaType),
			keyPage:    page,
			keyPerPage: perPage,
		})
		if err != nil {
			return err
		}

		return printMedia(cmd.OutOrStdout(), result.Media, g.asJSON, g.lang, g.noColor, g.mediaType)
	},
}

func init() {
	alltimeCmd.Flags().Int(keyPage, 1, "page number")
	alltimeCmd.Flags().Int(flagPerPage, 20, "results per page (max 50)")
	rootCmd.AddCommand(alltimeCmd)
}
