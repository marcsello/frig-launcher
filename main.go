package main

import (
	"log"
	"os"
	"os/signal"

	"syscall"

	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/bin/binttf"
	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
	"github.com/marcsello/frig-launcher/asset_management/fonts"
	"github.com/marcsello/frig-launcher/asset_management/image"
	"github.com/marcsello/frig-launcher/asset_management/loader"
	"github.com/marcsello/frig-launcher/asset_management/sound"
)

const (
	AssetStagePrimary = iota
	AssetStageSecondary
)

const (
	SNDLogo = iota
	SNDNavigate
	SNDSelect

	IMGTRex
	IMGFrig

	FontTitle
)

type FeedbackControllerWrapper struct{}

func (f *FeedbackControllerWrapper) OnNavigate() {
	_ = sound.PlaySnd(SNDNavigate)
}

func (f *FeedbackControllerWrapper) OnSelect() {
	_ = sound.PlaySnd(SNDSelect)
}

func (f *FeedbackControllerWrapper) OnError() {

}

func main() {
	defer binsdl.Load().Unload()
	defer binimg.Load().Unload()
	defer binttf.Load().Unload()
	defer sdl.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_GAMEPAD); err != nil {
		panic(err)
	}

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	err := sdl.DisableScreenSaver()
	if err != nil {
		log.Println("WARNING: Failed to inhibit screensaver: ", err)
	}

	displayID := sdl.GetPrimaryDisplay()
	log.Println("Primary display id: ", displayID)

	var displayName string
	displayName, err = displayID.Name()
	if err != nil {
		log.Println("WARNING: Failed to get primary display name", err)
	} else {
		log.Println("Primary display name: ", displayName)
	}

	dm, err := displayID.CurrentDisplayMode()
	if err != nil {
		panic(err)
	}

	log.Printf("Desktop: %dx%d @%fHz", dm.W, dm.H, dm.RefreshRate)

	window, renderer, err := sdl.CreateWindowAndRenderer("FRIG Launcher", int(dm.W), int(dm.H),
		sdl.WINDOW_FULLSCREEN|sdl.WINDOW_BORDERLESS|sdl.WINDOW_INPUT_FOCUS|sdl.WINDOW_ALWAYS_ON_TOP,
	)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	defer window.Destroy()

	err = sdl.HideCursor()
	if err != nil {
		log.Println("WARNING: Failed to hide cursor:", err)
	}

	var windowW, windowH int32
	windowW, windowH, err = window.Size()
	if err != nil {
		panic(err)
	}

	log.Printf("Window size: %dx%d", windowW, windowH)

	if sdl.HasGamepad() {
		log.Println("Gamepads detected")
	} else {
		log.Println("WARNING: No gamepads detected")
	}

	err = sound.Init()
	if err != nil {
		log.Println("WARNING: Audio init failed:", err)
	}
	defer sound.Close()

	image.Init(renderer)
	defer image.Close()

	err = fonts.Init(renderer)
	if err != nil {
		log.Println("WARNING: Fonts init failed:", err)
	}
	defer fonts.Close()

	// Assets needed at the very beginning
	loader.MustRegisterAsset(loader.SoundAsset, AssetStagePrimary, SNDLogo, "snd/logo.wav")
	loader.MustRegisterAsset(loader.ImageAsset, AssetStagePrimary, IMGTRex, "img/trex.png")

	// will be loaded while displaying the logo
	loader.MustRegisterAsset(loader.SoundAsset, AssetStageSecondary, SNDNavigate, "snd/navigate.wav")
	loader.MustRegisterAsset(loader.SoundAsset, AssetStageSecondary, SNDSelect, "snd/select.wav")
	loader.MustRegisterAsset(loader.ImageAsset, AssetStageSecondary, IMGFrig, "img/frig.png")
	loader.MustRegisterAsset(loader.FontAsset, AssetStageSecondary, FontTitle, "LiberationSans-Regular.ttf", loader.FontSize(32))

	loader.LoadAssetsNow(AssetStagePrimary)

	// play logo immediately
	err = sound.PlaySnd(SNDLogo)
	if err != nil {
		log.Println("Failed to play audio: ", err)
	}

	f := &FeedbackControllerWrapper{}

	sc := &SceneControllerImpl{}
	sc.RegisterNextScene(&LogoScene{})

	signalCh := make(chan os.Signal)
	go func() {
		sig := <-signalCh
		log.Println("Signal received: ", sig)
		sc.RequestExit()
	}()

	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	log.Println("start")
	err = sc.Run(renderer, f, int(windowW), int(windowH))
	if err != nil {
		panic(err)
	}
	log.Println("exit")
}
