package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/pavelsimo/ani/internal/display"
)

func printMedia(w io.Writer, media []anilist.Media, asJSON bool, lang string, noColor bool, mediaType string) error {
	if asJSON {
		return json.NewEncoder(w).Encode(media)
	}
	_, err := fmt.Fprintln(w, display.Render(media, lang, noColor, mediaType))
	return err
}
