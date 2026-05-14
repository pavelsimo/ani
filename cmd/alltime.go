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

var alltimeCmd = &cobra.Command{
	Use:          "all-time",
	Short:        "Show most popular anime of all time",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QueryAllTime, map[string]any{
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
	alltimeCmd.Flags().Int("page", 1, "page number")
	alltimeCmd.Flags().Int("per-page", 20, "results per page (max 50)")
	rootCmd.AddCommand(alltimeCmd)
}
