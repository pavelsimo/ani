package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/pavelsimo/ani/internal/display"
)

func printMedia(media []anilist.Media, asJSON bool, lang string, noColor bool, mediaType string) error {
	if asJSON {
		return json.NewEncoder(os.Stdout).Encode(media)
	}
	fmt.Println(display.Render(media, lang, noColor, mediaType))
	return nil
}
