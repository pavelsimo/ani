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

var topCmd = &cobra.Command{
	Use:          "top",
	Short:        "Show top anime by score",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		page, _ := cmd.Flags().GetInt("page")
		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")

		perPage := limit
		if perPage > 50 {
			perPage = 50
		}

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QueryTop, map[string]any{
			"page":    page,
			"perPage": perPage,
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
	topCmd.Flags().Int("limit", 20, "number of results (max 50 per page)")
	topCmd.Flags().Int("page", 1, "page number")
	rootCmd.AddCommand(topCmd)
}
