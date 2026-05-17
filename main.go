package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Zyko0/go-sdl3/bin/binimg"
	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/bin/binttf"
	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
	"github.com/marcsello/frig-launcher/asset_management/fonts"
	"github.com/marcsello/frig-launcher/asset_management/image"
	"github.com/marcsello/frig-launcher/asset_management/loader"
	"github.com/marcsello/frig-launcher/asset_management/sound"
	"github.com/marcsello/frig-launcher/config"
	"github.com/marcsello/frig-launcher/network"
	"github.com/marcsello/frig-launcher/utils"
	"gitlab.com/MikeTTh/env"
)

const (
	AssetStagePrimary = iota
	AssetStageSecondary
)

const (
	SNDLogo = iota
	SNDNavigate
	SNDSelect
)

const (
	IMGTRex = iota
	IMGFrig
	IMGNoNetwork
	IMGYesNetwork

	IMGIcon // this is basically a marker
)

const (
	FontTitle = iota
	FontClock
)

const (
	EVENT_MINUTE_PASSED = iota + sdl.EVENT_USER // it is not registered, but it does not seem necessary: https://github.com/libsdl-org/SDL/blob/a34b35a382bffcd430465e273264021fb0f2fab1/src/events/SDL_events.c#L1963-L1974
	EVENT_NETWORK_CHANGED
)

var (
	networkOnline atomic.Bool
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

type MuteFeedbackController struct{}

func (m MuteFeedbackController) OnNavigate() {}

func (m MuteFeedbackController) OnSelect() {}

func (m MuteFeedbackController) OnError() {}

func main() {
	var argShowSplash bool
	var argAudio bool
	var argWindow string
	var desiredFPS uint64
	flag.BoolVar(&argShowSplash, "splash", true, "Show splash screen")
	flag.BoolVar(&argAudio, "audio", true, "Enable audio")
	flag.Uint64Var(&desiredFPS, "desiredFPS", 30, "Limit FPS to desired value")
	flag.StringVar(&argWindow, "window", "", "Use windowed mode with the specified resolution. Example -window=1024x768 otherwise will be fullscreen")
	flag.Parse()

	defer binsdl.Load().Unload()
	defer binimg.Load().Unload()
	defer binttf.Load().Unload()
	defer sdl.Quit()

	var err error

	err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO | sdl.INIT_GAMEPAD)
	if err != nil {
		panic(err)
	}

	err = ttf.Init()
	if err != nil {
		panic(err)
	}

	if !env.Bool("FRIG_NO_SCREENSAVER_INHIBIT", false) {
		err = sdl.DisableScreenSaver()
		if err != nil {
			log.Println("WARNING: Failed to inhibit screensaver: ", err)
		}
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

	var windowFlags sdl.WindowFlags
	var windowW, windowH int32

	if argWindow == "" {
		windowFlags = sdl.WINDOW_FULLSCREEN | sdl.WINDOW_BORDERLESS | sdl.WINDOW_INPUT_FOCUS | sdl.WINDOW_ALWAYS_ON_TOP
		windowW = dm.W
		windowH = dm.H
	} else {
		windowW, windowH, err = utils.ParseResolution(argWindow)
		if err != nil {
			panic(err)
		}
	}

	window, renderer, err := sdl.CreateWindowAndRenderer("FRIG Launcher", int(windowW), int(windowH),
		windowFlags,
	)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	defer window.Destroy()

	// Disable some unneeded events
	sdl.SetEventEnabled(sdl.EVENT_TEXT_INPUT, false)
	sdl.SetEventEnabled(sdl.EVENT_SENSOR_UPDATE, false)
	sdl.SetEventEnabled(sdl.EVENT_MOUSE_MOTION, false)
	sdl.SetEventEnabled(sdl.EVENT_WINDOW_MOUSE_ENTER, false)
	sdl.SetEventEnabled(sdl.EVENT_WINDOW_MOUSE_LEAVE, false)
	sdl.SetEventEnabled(sdl.EVENT_MOUSE_BUTTON_UP, false)
	sdl.SetEventEnabled(sdl.EVENT_MOUSE_BUTTON_DOWN, false)
	sdl.SetEventEnabled(sdl.EVENT_MOUSE_WHEEL, false)
	sdl.SetEventEnabled(sdl.EVENT_WINDOW_MOVED, false)
	sdl.SetEventEnabled(sdl.EVENT_CLIPBOARD_UPDATE, false)

	err = sdl.HideCursor()
	if err != nil {
		log.Println("WARNING: Failed to hide cursor:", err)
	}

	// re-read the value to "verify" it
	windowW, windowH, err = window.Size()
	if err != nil {
		panic(err)
	}

	log.Printf("Window size: %dx%d", windowW, windowH)

	if sdl.HasGamepad() {
		log.Println("Gamepads detected")
	} else {
		log.Println("WARNING: No gamepads detected!")
	}

	if argAudio {
		err = sound.Init()
		if err != nil {
			log.Println("WARNING: Audio init failed:", err)
		}
		defer sound.Close()
	}

	image.Init(renderer)
	defer image.Close()

	err = fonts.Init(renderer)
	if err != nil {
		log.Println("WARNING: Fonts init failed:", err)
	}
	defer fonts.Close()

	// Setup interrupt handler
	signalCh := make(chan os.Signal)
	go func() {
		sig := <-signalCh
		log.Println("Signal received:", sig)
		err := sdl.PushEvent(&sdl.Event{Type: sdl.EVENT_QUIT}) // the docs state: It is safe to call this function from any thread.
		if err != nil {
			log.Println("failed to push quit event")
		}
	}()

	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	// Setup clock ticker
	go func() {
		for {
			time.Sleep(time.Now().Add(time.Minute).Truncate(time.Minute).Sub(time.Now()))
			evt := sdl.UserEvent{
				Type: uint32(EVENT_MINUTE_PASSED), // I think the library is broken here...
			}
			err := utils.PushUserEvent(&evt)
			if err != nil {
				log.Println("Failed to push minute passed event:", err)
			}
		}
	}()

	// Setup network checker
	go func() {
		firstLoop := true
		for {
			state := network.HasNetwork()
			lastState := networkOnline.Swap(state)

			if state != lastState || firstLoop {
				log.Println("Network state changed to", state)
				ev := sdl.UserEvent{
					Type: uint32(EVENT_NETWORK_CHANGED),
					Code: int32(utils.BoolToInt(state)),
				}
				err = utils.PushUserEvent(&ev)
				if err != nil {
					log.Println("Failed to push network changed event:", err)
				}
			}

			firstLoop = false
			time.Sleep(time.Second * 10) // sleep 10 seconds before polling the network again...
		}
	}()

	// That's enough "framework" setup, now let's run the "application logic"

	err = config.LoadConfig()
	if err != nil {
		log.Println("Failed to load config")
		panic(err)
	}
	if len(config.Config.Applications) == 0 {
		log.Println("WARNING: No applications defined!!!")
	}

	// Assets needed at the very beginning
	if argAudio {
		loader.MustRegisterAsset(loader.SoundAsset, AssetStagePrimary, SNDLogo, "snd/logo.wav")
	}
	loader.MustRegisterAsset(loader.ImageAsset, AssetStagePrimary, IMGTRex, "img/trex.png")

	// will be loaded while displaying the logo
	if argAudio {
		loader.MustRegisterAsset(loader.SoundAsset, AssetStageSecondary, SNDNavigate, "snd/navigate.wav")
		loader.MustRegisterAsset(loader.SoundAsset, AssetStageSecondary, SNDSelect, "snd/select.wav")
	}
	loader.MustRegisterAsset(loader.ImageAsset, AssetStageSecondary, IMGFrig, "img/frig.png")
	loader.MustRegisterAsset(loader.ImageAsset, AssetStageSecondary, IMGNoNetwork, "img/no_network.svg", loader.SVGScaleByH(32)) // same scale as text
	loader.MustRegisterAsset(loader.ImageAsset, AssetStageSecondary, IMGYesNetwork, "img/yes_network.svg", loader.SVGScaleByH(32))

	for i, app := range config.Config.Applications {
		loader.MustRegisterAsset(loader.ImageAsset, AssetStageSecondary, IMGIcon+i, app.Icon, loader.SVGScaleByH(windowH/2)) // the largest icon size ever displayed
	}

	loader.MustRegisterAsset(loader.FontAsset, AssetStageSecondary, FontTitle, "NotoSans-Bold.ttf", loader.FontSize(48))
	loader.MustRegisterAsset(loader.FontAsset, AssetStageSecondary, FontClock, "NotoSans-Bold.ttf", loader.FontSize(32))

	loader.LoadAssetsNow(AssetStagePrimary)

	if argAudio && argShowSplash {
		// play logo immediately
		err = sound.PlaySnd(SNDLogo)
		if err != nil {
			log.Println("Failed to play audio: ", err)
		}
	}

	var f FeedbackController
	if argAudio {
		f = &FeedbackControllerWrapper{}
	} else {
		f = &MuteFeedbackController{}
	}

	sc := &SceneControllerImpl{DesiredFPS: desiredFPS}
	if argShowSplash {
		sc.RegisterNextScene(&LogoScene{})
	} else {
		loader.LoadAssetsNow(AssetStageSecondary)
		sc.RegisterNextScene(&MenuScene{})
	}

	log.Println("start")
	err = sc.Run(renderer, f, int(windowW), int(windowH))
	if err != nil {
		panic(err)
	}
	log.Println("clean exit")
}
