package loader

import (
	"fmt"
	"log"
	"path"

	"github.com/marcsello/frig-launcher/asset_management/fonts"
	"github.com/marcsello/frig-launcher/asset_management/image"
	"github.com/marcsello/frig-launcher/asset_management/sound"
	"github.com/marcsello/frig-launcher/utils"
)

type Asset struct {
	Type        AssetType
	Stage       int // user defined stage
	ID          int // user defined id used for accessing the resource by the asset subsystems
	AbsFilePath string

	FontSize  float32 // specific to font only
	SVGScaleW int32   // specific to img type of svg only
	SVGScaleH int32   // specific to img type of svg only

	Loaded bool
}

type ExtraConfig func(*Asset)

func FontSize(ptSize float32) ExtraConfig {
	return func(asset *Asset) {
		asset.FontSize = ptSize
	}
}

func SVGScale(w, h int32) ExtraConfig {
	return func(asset *Asset) {
		asset.SVGScaleW = w
		asset.SVGScaleH = h
	}
}

func SVGScaleByW(w int32) ExtraConfig {
	return func(asset *Asset) {
		asset.SVGScaleW = w
		asset.SVGScaleH = 0
	}
}

func SVGScaleByH(h int32) ExtraConfig {
	return func(asset *Asset) {
		asset.SVGScaleW = 0
		asset.SVGScaleH = h
	}
}

var (
	assets []Asset
)

// RegisterAsset takes a relative path, it will attempt to resolve the path for the asset
func RegisterAsset(t AssetType, stage, id int, assetPath string, extraConfig ...ExtraConfig) error {
	var absPath string
	if path.IsAbs(assetPath) {
		if utils.IsFile(assetPath) {
			absPath = assetPath
		}
	} else {
		absPath = resolvePath(t, assetPath)
	}

	if absPath == "" {
		log.Printf("Could not locate %s asset: %s Attempting to load it will fail!", t.String(), assetPath)
		return fmt.Errorf("could not locate %s asset", t.String())
	}
	log.Printf("Located %s asset: %s", t.String(), absPath)

	record := Asset{
		Type:        t,
		Stage:       stage,
		ID:          id,
		AbsFilePath: absPath,
		Loaded:      false,
	}

	for _, fn := range extraConfig {
		fn(&record)
	}

	assets = append(assets, record)
	return nil
}

func MustRegisterAsset(t AssetType, stage, id int, relPath string, extraConfig ...ExtraConfig) {
	err := RegisterAsset(t, stage, id, relPath, extraConfig...)
	if err != nil {
		log.Printf("FATAL: Failed to register %s asset! %v", t.String(), err)
		panic(err)
	}
}

func loadAsset(asset Asset) error {
	switch asset.Type {
	case SoundAsset:
		return sound.LoadSnd(asset.ID, asset.AbsFilePath)
	case ImageAsset:
		return image.LoadImageResource(asset.ID, asset.AbsFilePath, asset.SVGScaleW, asset.SVGScaleH)
	case FontAsset:
		return fonts.LoadFont(asset.ID, asset.AbsFilePath, asset.FontSize)
	default:
		return nil
	}
}

// LoadAssetsNow loads all assets required by `stage`
func LoadAssetsNow(stage int) {
	for _, asset := range assets {
		if asset.Stage != stage {
			continue
		}
		if asset.Loaded {
			log.Printf("Attempted to load %s asset twice: %s", asset.Type.String(), asset.AbsFilePath)
			continue
		}

		err := loadAsset(asset)
		if err != nil {
			log.Printf("Failed to load %s %s: %v", asset.Type.String(), asset.AbsFilePath, err)
			continue
		}

		log.Printf("%s asset loaded: %s to %d", asset.Type.String(), asset.AbsFilePath, asset.ID)
		asset.Loaded = true
	}
}
