package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pavelsimo/ani/internal/anilist"
)

// newClient builds the AniList client; tests swap it for a mocked one.
var newClient = anilist.New

// flagPerPage is the CLI flag name; keyPerPage is the GraphQL variable name.
const flagPerPage = "per-page"

// maxPerPage is the AniList Page(perPage:) hard cap.
const maxPerPage = 50

// clampPerPage clamps n into [1, maxPerPage].
func clampPerPage(n int) int {
	if n < 1 {
		return 1
	}
	if n > maxPerPage {
		return maxPerPage
	}
	return n
}

// pageFlags reads and validates the pagination flags. --page must be >= 1;
// per-page is clamped to [1, maxPerPage]. Commands without a per-page flag
// get perPage 0.
func pageFlags(cmd *cobra.Command) (page, perPage int, err error) {
	page, _ = cmd.Flags().GetInt(keyPage)
	if page < 1 {
		return 0, 0, fmt.Errorf("--page must be >= 1")
	}
	if cmd.Flags().Lookup(flagPerPage) != nil {
		perPage, _ = cmd.Flags().GetInt(flagPerPage)
		perPage = clampPerPage(perPage)
	}
	return page, perPage, nil
}

// globalFlags holds the persistent flags every command reads.
type globalFlags struct {
	asJSON    bool
	noColor   bool
	lang      string
	mediaType string
}

// getGlobalFlags reads the persistent flags. Retrieval errors are ignored:
// they only occur for unregistered names or type mismatches, which are
// programmer errors, not runtime conditions.
func getGlobalFlags(cmd *cobra.Command) globalFlags {
	asJSON, _ := cmd.Flags().GetBool("json")
	noColor, _ := cmd.Flags().GetBool("no-color")
	lang, _ := cmd.Flags().GetString("lang")
	mediaType, _ := cmd.Flags().GetString(keyType)
	return globalFlags{asJSON: asJSON, noColor: noColor, lang: lang, mediaType: mediaType}
}
