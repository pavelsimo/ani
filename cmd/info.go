package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

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
			return fmt.Errorf("AniList ID must be a positive integer, got %q", args[0])
		}

		g := getGlobalFlags(cmd)

		client := newClient()
		media, err := client.QueryMedia(context.Background(), id)
		if err != nil {
			return err
		}

		if g.asJSON {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(media)
		}
		_, err = fmt.Fprint(cmd.OutOrStdout(), display.RenderDetailWithOptions(*media, g.lang, display.DetailOptions{
			Width:     80,
			NoColor:   g.noColor,
			MediaType: media.Type,
		}))
		return err
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
