package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
)

var popularCmd = &cobra.Command{
	Use:          "popular",
	Short:        "Show popular anime for the current or specified season",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		season, _ := cmd.Flags().GetString("season")
		year, _ := cmd.Flags().GetInt("year")
		page, _ := cmd.Flags().GetInt(keyPage)
		perPage, _ := cmd.Flags().GetInt("per-page")
		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")
		mediaType, _ := cmd.Flags().GetString(keyType)

		vars := map[string]any{
			keyType:    strings.ToUpper(mediaType),
			keyPage:    page,
			keyPerPage: perPage,
		}

		// Season filtering only applies to anime.
		if strings.ToUpper(mediaType) != "MANGA" {
			if season == "" || year == 0 {
				s, y := anilist.CurrentSeason()
				if season == "" {
					season = s
				}
				if year == 0 {
					year = y
				}
			}
			vars["season"] = strings.ToUpper(season)
			vars["seasonYear"] = year
		}

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QueryPopularSeason, vars)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return err
		}

		return printMedia(result.Media, asJSON, lang, noColor, mediaType)
	},
}

func init() {
	popularCmd.Flags().String("season", "", "season: winter, spring, summer, fall (default: current)")
	popularCmd.Flags().Int("year", 0, "year (default: current year)")
	popularCmd.Flags().Int(keyPage, 1, "page number")
	popularCmd.Flags().Int("per-page", 20, "results per page (max 50)")
	rootCmd.AddCommand(popularCmd)
}
