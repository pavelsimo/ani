package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/pavelsimo/ani/internal/display"
)

var popularCmd = &cobra.Command{
	Use:          "popular",
	Short:        "Show popular anime for the current or specified season",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		season, _ := cmd.Flags().GetString("season")
		year, _ := cmd.Flags().GetInt("year")
		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")

		if season == "" || year == 0 {
			s, y := anilist.CurrentSeason()
			if season == "" {
				season = s
			}
			if year == 0 {
				year = y
			}
		}

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QueryPopularSeason, map[string]any{
			"season":     strings.ToUpper(season),
			"seasonYear": year,
			"page":       page,
			"perPage":    perPage,
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
	popularCmd.Flags().String("season", "", "season: winter, spring, summer, fall (default: current)")
	popularCmd.Flags().Int("year", 0, "year (default: current year)")
	popularCmd.Flags().Int("page", 1, "page number")
	popularCmd.Flags().Int("per-page", 20, "results per page (max 50)")
	rootCmd.AddCommand(popularCmd)
}
