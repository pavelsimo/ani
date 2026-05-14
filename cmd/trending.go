package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/pavelsimo/ani/internal/display"
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

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QueryTrending, map[string]any{
			keyPage:    page,
			keyPerPage: perPage,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return err
		}

		if asJSON {
			return json.NewEncoder(os.Stdout).Encode(result.Media)
		}
		fmt.Print(display.Render(result.Media, lang, noColor))
		return nil
	},
}

func init() {
	trendingCmd.Flags().Int(keyPage, 1, "page number")
	trendingCmd.Flags().Int("per-page", 20, "results per page (max 50)")
	rootCmd.AddCommand(trendingCmd)
}
