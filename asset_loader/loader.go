package asset_loader

import (
	"log"
	"os"
	"path"

	"github.com/marcsello/frig-launcher/image"
	"github.com/marcsello/frig-launcher/sound"
	"github.com/marcsello/frig-launcher/utils"
	"gitlab.com/MikeTTh/env"
)

var globalAssetLocations = []string{ // in order of preference
	"/usr/share/frig/assets",
	"/usr/local/share/games/frig/assets",
	"/opt/frig/assets",
}

var (
	soundAssets = make(map[int]map[int]string)
	imageAssets = make(map[int]map[int]string)

	stagesAlreadyLoaded = make(map[int]struct{})
)

func resolvePath(relPath string) string {
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

	// give up
	return ""
}

func RegisterSoundAsset(stage, id int, relPath string) {
	if soundAssets[stage] == nil {
		soundAssets[stage] = make(map[int]string)
	}

	absPath := resolvePath(relPath)
	if absPath == "" {
		log.Printf("Could not locate sound asset: %s Attempting to use it will fail!", relPath)
		return
	}
	log.Printf("Located sound asset: %s", absPath)

	soundAssets[stage][id] = absPath
}

func RegisterImageAsset(stage, id int, relPath string) {
	if imageAssets[stage] == nil {
		imageAssets[stage] = make(map[int]string)
	}

	absPath := resolvePath(relPath)
	if absPath == "" {
		log.Printf("Could not locate image asset: %s Attempting to use it will fail!", relPath)
		return
	}
	log.Printf("Located image asset: %s", absPath)

	imageAssets[stage][id] = absPath
}

func loadAssets(sounds, images map[int]string) {
	for id, absPath := range sounds {
		err := sound.LoadSnd(id, absPath)
		if err != nil {
			log.Printf("Failed to load sound %s: %v", absPath, err)
		} else {
			log.Printf("Sound asset loaded: %s to %d", absPath, id)
		}
	}
	for id, absPath := range images {
		err := image.LoadImageResource(id, absPath)
		if err != nil {
			log.Printf("Failed to load image %s: %v", absPath, err)
		} else {
			log.Printf("Image asset loaded: %s to %d", absPath, id)
		}
	}
}

// LoadAssetsNow loads all assets required by `stage`
func LoadAssetsNow(stage int) {
	_, ok := stagesAlreadyLoaded[stage]
	if ok {
		log.Printf("Attempted to load stage %d twice!", stage)
		return
	}

	loadAssets(soundAssets[stage], imageAssets[stage])
	stagesAlreadyLoaded[stage] = struct{}{}
}
