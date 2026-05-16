package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/pavelsimo/ani/internal/display"
)

var infoCmd = &cobra.Command{
	Use:          "info <id>",
	Short:        "Show full details for an anime by AniList ID",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil || id <= 0 {
			fmt.Fprintln(os.Stderr, "error: AniList ID must be a positive integer")
			return fmt.Errorf("invalid AniList ID: %s", args[0])
		}

		asJSON, _ := cmd.Flags().GetBool("json")
		noColor, _ := cmd.Flags().GetBool("no-color")
		lang, _ := cmd.Flags().GetString("lang")

		client := anilist.New()
		media, err := client.QueryMedia(context.Background(), id)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return err
		}

		if asJSON {
			return json.NewEncoder(os.Stdout).Encode(media)
		}
		fmt.Print(display.RenderDetailWithOptions(*media, lang, display.DetailOptions{
			Width:     80,
			NoColor:   noColor,
			MediaType: media.Type,
		}))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
