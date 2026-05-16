package cmd

import (
	"context"
	"fmt"
	"os"
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
		page, _ := cmd.Flags().GetInt(keyPage)
		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")
		mediaType, _ := cmd.Flags().GetString("type")

		perPage := limit
		if perPage > 50 {
			perPage = 50
		}

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QueryTop, map[string]any{
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
	topCmd.Flags().Int("limit", 20, "number of results (max 50 per page)")
	topCmd.Flags().Int(keyPage, 1, "page number")
	rootCmd.AddCommand(topCmd)
}
