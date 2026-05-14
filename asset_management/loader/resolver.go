package loader

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/liamg/fontinfo"
	"github.com/marcsello/frig-launcher/utils"
	"gitlab.com/MikeTTh/env"
)

var (
	fontsList            []fontinfo.Font
	globalAssetLocations = []string{ // in order of preference
		"/usr/share/frig/assets",
		"/usr/local/share/games/frig/assets",
		"/opt/frig/assets",
	}
)

func searchSystemFont(filename string) string {
	if fontsList == nil {
		var err error
		fontsList, err = fontinfo.List()
		if err != nil {
			log.Println("WARNING: Failed to list system fonts:", err, " System fonts will not be loadable!")
			return ""
		}
	}

	for _, font := range fontsList {
		if strings.EqualFold(path.Base(font.Path), filename) {
			return font.Path
		}
	}

	// Maybe it is the family name of the font???
	for _, font := range fontsList {
		if strings.EqualFold(font.Family, filename) {
			return font.Path
		}
	}

	// give up
	return ""
}

func resolvePath(assetType AssetType, relPath string) string {
	override := env.String("FRIG_ASSETS_PATH", "")
	if override != "" {
		absPath := path.Join(override, relPath)
		if utils.IsFile(absPath) {
			return absPath
		}
		return "" // If overridden, then don't attempt the regular lookup...
	}

	// attempt the regular global paths
	for _, rootPath := range globalAssetLocations {
		absPath := path.Join(rootPath, relPath)
		if utils.IsFile(absPath) {
			return absPath
		}
	}

	if !env.Bool("FRIG_ASSETS_NOCWD", false) { // maybe this has some security implications
		// last resort: try cwd
		cwd, err := os.Getwd()
		if err != nil {
			log.Println("Failed to acquire CWD:", err)
			return ""
		}

		absPath := path.Join(cwd, "assets", relPath)
		if utils.IsFile(absPath) {
			return absPath
		}
	}

	if assetType == FontAsset {
		// maybe we are looking for a system font?
		absPath := searchSystemFont(relPath)
		if absPath != "" {
			return absPath
		}
	}

	// give up
	return ""
}
