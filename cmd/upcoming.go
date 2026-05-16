package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
)

var upcomingCmd = &cobra.Command{
	Use:          "upcoming",
	Short:        "Show upcoming anime not yet airing",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt(keyPage)
		perPage, _ := cmd.Flags().GetInt("per-page")
		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")
		mediaType, _ := cmd.Flags().GetString(keyType)

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QueryUpcoming, map[string]any{
			keyType:    strings.ToUpper(mediaType),
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
	upcomingCmd.Flags().Int(keyPage, 1, "page number")
	upcomingCmd.Flags().Int("per-page", 20, "results per page (max 50)")
	rootCmd.AddCommand(upcomingCmd)
}
