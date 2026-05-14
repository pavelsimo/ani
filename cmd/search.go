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

var searchCmd = &cobra.Command{
	Use:          "search [query]",
	Short:        "Search anime by title, genre, year, season, or format",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		genres, _ := cmd.Flags().GetStringArray("genre")
		year, _ := cmd.Flags().GetInt("year")
		season, _ := cmd.Flags().GetString("season")
		format, _ := cmd.Flags().GetString("format")
		status, _ := cmd.Flags().GetString("status")
		minScore, _ := cmd.Flags().GetInt("min-score")
		page, _ := cmd.Flags().GetInt(keyPage)
		perPage, _ := cmd.Flags().GetInt("per-page")
		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")

		vars := map[string]any{
			keyPage:    page,
			keyPerPage: perPage,
		}

		if len(args) > 0 {
			vars["search"] = strings.Join(args, " ")
		}
		if len(genres) > 0 {
			vars["genres"] = genres
		}
		if year > 0 {
			vars["seasonYear"] = year
		}
		if season != "" {
			vars["season"] = strings.ToUpper(season)
		}
		if format != "" {
			vars["format"] = strings.ToUpper(format)
		}
		if status != "" {
			vars["status"] = statusEnum(status)
		}
		if minScore > 0 {
			vars["averageScore_greater"] = minScore
		}

		if len(vars) <= 2 {
			fmt.Fprintln(os.Stderr, "error: provide at least a search query or one filter flag")
			return fmt.Errorf("no search criteria provided")
		}

		client := anilist.New()
		result, err := client.Query(context.Background(), anilist.QuerySearch, vars)
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

func statusEnum(s string) string {
	switch strings.ToLower(s) {
	case "airing", "releasing":
		return "RELEASING"
	case "finished":
		return "FINISHED"
	case "upcoming", "not_yet_released":
		return "NOT_YET_RELEASED"
	case "cancelled":
		return "CANCELLED"
	default:
		return strings.ToUpper(s)
	}
}

func init() {
	searchCmd.Flags().StringArray("genre", nil, "filter by genre (repeatable)")
	searchCmd.Flags().Int("year", 0, "filter by year")
	searchCmd.Flags().String("season", "", "season: winter, spring, summer, fall")
	searchCmd.Flags().String("format", "", "format: tv, ova, ona, movie, special")
	searchCmd.Flags().String("status", "", "status: airing, finished, upcoming")
	searchCmd.Flags().Int("min-score", 0, "minimum average score (0-100)")
	searchCmd.Flags().Int(keyPage, 1, "page number")
	searchCmd.Flags().Int("per-page", 20, "results per page (max 50)")
	rootCmd.AddCommand(searchCmd)
}
