package cmd

import (
	"context"
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
		page, perPage, err := pageFlags(cmd)
		if err != nil {
			return err
		}
		g := getGlobalFlags(cmd)

		vars := map[string]any{
			keyType:    strings.ToUpper(g.mediaType),
			keyPage:    page,
			keyPerPage: perPage,
		}

		// Season filtering only applies to anime.
		if strings.ToUpper(g.mediaType) != "MANGA" {
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

		client := newClient()
		result, err := client.Query(context.Background(), anilist.QueryPopularSeason, vars)
		if err != nil {
			return err
		}

		return printMedia(cmd.OutOrStdout(), result.Media, g.asJSON, g.lang, g.noColor, g.mediaType)
	},
}

func init() {
	popularCmd.Flags().String("season", "", "season: winter, spring, summer, fall (default: current)")
	popularCmd.Flags().Int("year", 0, "year (default: current year)")
	popularCmd.Flags().Int(keyPage, 1, "page number")
	popularCmd.Flags().Int(flagPerPage, 20, "results per page (max 50)")
	rootCmd.AddCommand(popularCmd)
}
