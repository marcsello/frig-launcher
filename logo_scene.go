package main

import (
	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/marcsello/frig-launcher/asset_management/image"
	"github.com/marcsello/frig-launcher/asset_management/loader"
	"github.com/marcsello/frig-launcher/utils"
)

type LogoScene struct {
	fbc FeedbackController
	sc  SceneController

	logoTexture *sdl.Texture

	finished     bool
	assetsLoaded bool
}

func (s *LogoScene) Draw(renderer *sdl.Renderer, scrW, scrH int, firstFrame bool, dtNS, durationNS uint64) (bool, error) {

	trex, _ := image.GetTexture(IMGTRex)

	durationMs := float64(durationNS) / 1000.0 / 1000.0

	x := scrW / 2
	y := scrH / 2

	var value uint8
	var frigWScale float64
	var xShift float64

	switch {
	case durationMs < 250:
		value = uint8(255 * (durationMs / 250))
	case durationMs > 4000 && durationMs < 5000:
		value = 255 - uint8(255*((durationMs-4000)/1000))
	case durationMs > 5000:
		value = 0
	default:
		value = 255
	}

	if durationMs > 1500 && durationMs < 2000 {
		frigWScale = utils.EaseInOutCubic((durationMs - 1500) / 500)
		xShift = 150 * frigWScale
	} else if durationMs >= 2000 {
		frigWScale = 1
		xShift = 150
	}

	if durationMs > 5500 {
		s.finished = true
	}

	if durationMs > 250 && !s.assetsLoaded { // load assets after displaying the logo, but before unrolling the text
		loader.LoadAssetsNow(AssetStageSecondary)
		s.assetsLoaded = true
	}

	err := renderer.SetDrawColor(0, 0, 0, 255)
	if err != nil {
		return false, err
	}
	err = renderer.Clear()
	if err != nil {
		return false, err
	}
	err = trex.SetAlphaMod(value)
	if err != nil {
		return false, err
	}

	err = renderer.RenderTexture(trex, nil, &sdl.FRect{
		X: float32(float64(x-50) + xShift),
		Y: float32(y - 50),
		W: 91, // (100/139)*127
		H: 100,
	})
	if err != nil {
		return false, err
	}

	if durationMs > 1500 && s.assetsLoaded {
		frig, _ := image.GetTexture(IMGFrig)

		err = frig.SetAlphaMod(value)
		if err != nil {
			return false, err
		}
		err = renderer.RenderTexture(frig, &sdl.FRect{
			X: 0,
			Y: 0,
			W: float32(300 * frigWScale),
			H: 100,
		}, &sdl.FRect{
			X: float32(float64(x-55) - xShift),
			Y: float32(y - 50),
			W: float32(300 * frigWScale),
			H: 100,
		})
		if err != nil {
			return false, err
		}
	}

	return !s.finished, nil
}

func (s *LogoScene) MaxInhibitMS() int32 {
	return 0 // don't allow inhibition
}

func (s *LogoScene) Bind(fbc FeedbackController, sc SceneController) {
	s.fbc = fbc
	s.sc = sc

	sc.RegisterNextScene(&MenuScene{})
}

func (s *LogoScene) UnBind() {
	if !s.assetsLoaded {
		loader.LoadAssetsNow(AssetStageSecondary)
		s.assetsLoaded = true
	}

	s.fbc = nil
	s.sc = nil
}

func (s *LogoScene) Input(_ UserIntent) {
	// The user wants this to be skipped
	s.finished = true
}

func (s *LogoScene) WallClockMinutePassed() {} // noop for logo
