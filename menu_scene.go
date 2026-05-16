package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/marcsello/frig-launcher/asset_management/fonts"
	"github.com/marcsello/frig-launcher/asset_management/image"
	"github.com/marcsello/frig-launcher/config"
	"github.com/marcsello/frig-launcher/utils"
)

type MenuScene struct {
	fbc FeedbackController
	sc  SceneController

	selection     int
	lastSelection int

	choiceMade   bool
	launched     bool
	choiceMadeAt uint64

	iconsAnimators []*utils.TransitionAnimator
}

func (m *MenuScene) Draw(renderer *sdl.Renderer, scrW, scrH int, firstFrame bool, dtNS, durationNS uint64) (bool, error) {
	var err error

	var anythingMoved bool

	inactiveIconSize := float32(scrW / 6)
	activeIconSize := float32(scrH / 2)

	if firstFrame {
		// Initialize these with random off-screen positions
		radius := math.Max(float64(scrW), float64(scrH)) + float64(activeIconSize)
		for i := range m.iconsAnimators {
			rad := rand.Float64() * 2 * math.Pi
			x := float64(scrW/2) + radius*math.Cos(rad)
			y := float64(scrH/2) + radius*math.Sin(rad)

			m.iconsAnimators[i] = &utils.TransitionAnimator{Rect: sdl.FRect{
				X: float32(x),
				Y: float32(y),
				W: inactiveIconSize,
				H: inactiveIconSize,
			}}
		}
		anythingMoved = true // First frame must be drawn
	}

	if m.choiceMade && m.choiceMadeAt == 0 {
		m.choiceMadeAt = durationNS
	}

	var unselectedValue, selectedValue, textValue uint8
	if m.choiceMade {
		since := durationNS - m.choiceMadeAt

		if since < 100_000_000 {
			unselectedValue = 255 - uint8(255*(float64(since)/100_000_000)) // 100ms fade-out
			textValue = unselectedValue
		} else {
			unselectedValue = 0
			textValue = 0
		}

		if since < 200_000_000 {
			selectedValue = 255
		} else if since < 550_000_000 {
			selectedValue = 255 - uint8(255*(float64(since-350_000_000)/350_000_000))
		} else {
			selectedValue = 0
		}

		if since > 550_000_000 && !m.launched {
			log.Println("launch")
			Launch(config.Config.Applications[m.selection].Exec)
			m.launched = true
		}

		if since < 1_000_000_000 { // keep the empty screen for 1 sec, to allow the sound effect to play properly
			anythingMoved = true
		} else {
			// finally exit
			return false, sdl.EndLoop
		}

	} else if durationNS < 1_000_000_000 {
		unselectedValue = 255
		selectedValue = 255
		textValue = uint8(255 * (float64(durationNS) / 1_000_000_000)) // 1 sec fade-in
		anythingMoved = true
	} else {
		unselectedValue = 255
		selectedValue = 255
		textValue = 255
	}

	if m.lastSelection != m.selection {
		for i := range m.iconsAnimators {

			x := (float32(scrW) - activeIconSize) / 2 // center by default

			xShift := i - m.selection // negative when standing before, positive if after

			if xShift > 0 {
				// this is true black magic :D
				// And I'm not explaining it
				// if it was hard to come up with, it should be hard to debug
				// ...
				// Just kidding, actually this just adds an extra padding if the icons standing AFTER the selected icon
				// This is needed because, we are calculating the left side of every icon
				x += activeIconSize - inactiveIconSize
			}
			x += inactiveIconSize * float32(xShift)

			if xShift != 0 {
				x += (inactiveIconSize / 8) * float32(xShift)
			}

			var rect sdl.FRect
			if i == m.selection {
				rect = sdl.FRect{
					X: x,
					Y: (float32(scrH) - activeIconSize) / 2,
					W: activeIconSize,
					H: activeIconSize,
				}
			} else {
				rect = sdl.FRect{
					X: x,
					Y: (float32(scrH) / 2) - (inactiveIconSize / 4),
					W: inactiveIconSize,
					H: inactiveIconSize,
				}
			}

			m.iconsAnimators[i].SetDesired(rect)
		}
		m.lastSelection = m.selection
		anythingMoved = true // selection changed, something must have moved
	}

	for i := range m.iconsAnimators {
		if m.iconsAnimators[i].Update(dtNS) {
			anythingMoved = true
		}
	}

	// This was all the calculations basically, so now we just draw everything

	err = renderer.SetDrawColor(0, 0, 0, 255)
	if err != nil {
		return false, err
	}

	err = renderer.Clear()
	if err != nil {
		return false, err
	}

	for i := range config.Config.Applications {
		var v uint8
		if i == m.selection {
			v = selectedValue
		} else {
			v = unselectedValue
		}

		var icon *sdl.Texture
		icon, _ = image.GetTexture(IMGIcon + i)
		err = icon.SetAlphaMod(v)
		if err != nil {
			return false, err
		}

		anim := m.iconsAnimators[i]
		if anim.Rect.X > float32(scrW) || anim.Rect.X < -anim.Rect.W || anim.Rect.Y > float32(scrH) || anim.Rect.Y < -anim.Rect.H {
			// no need to render these
			continue
		}

		err = renderer.RenderTexture(icon, nil, &anim.Rect)
		if err != nil {
			return false, err
		}
	}

	textTitle, err := fonts.GetText(FontTitle, config.Config.Applications[m.selection].Name)
	if err != nil {
		return false, err
	}

	textClock, err := fonts.GetText(FontClock, "12:23")
	if err != nil {
		return false, err
	}

	err = textTitle.SetColor(255, 255, 255, textValue)
	if err != nil {
		return false, err
	}

	err = textClock.SetColor(255, 255, 255, textValue)
	if err != nil {
		return false, err
	}

	textTitleW, _, err := textTitle.Size()
	if err != nil {
		return false, err
	}

	err = textTitle.DrawRenderer(float32(scrW/2)-float32(textTitleW/2), float32(scrH-150))
	if err != nil {
		return false, err
	}

	textClockW, _, err := textClock.Size()
	if err != nil {
		return false, err
	}

	err = textClock.DrawRenderer(float32(scrW-32-int(textClockW)), 32)
	if err != nil {
		return false, err
	}

	return anythingMoved, nil
}

func (m *MenuScene) MaxInhibitMS() int32 {
	return 15_000 // unfreeze every 15 sec, to keep the clock ticking
}

func (m *MenuScene) Bind(fbc FeedbackController, sc SceneController) {
	m.fbc = fbc
	m.sc = sc

	m.iconsAnimators = make([]*utils.TransitionAnimator, len(config.Config.Applications)) // this will be all nil at first

	// try to select the icon at the middle
	m.selection = len(config.Config.Applications)/2 - 1
	if m.selection < 0 {
		m.selection = 0
	}
	m.lastSelection = -1
}

func (m *MenuScene) UnBind() {
	m.fbc = nil
	m.sc = nil
}

func (m *MenuScene) Input(intent UserIntent) {
	if m.choiceMade {
		// input refused
		return
	}

	//goland:noinspection GoSwitchMissingCasesForIotaConsts
	switch intent {
	case IntentNext:
		m.inputNext()
	case IntentPrev:
		m.inputPrev()
	case IntentSelect:
		m.inputSelect()
	}
}

func (m *MenuScene) inputNext() {
	if m.selection == len(config.Config.Applications)-1 {
		return
	}

	m.selection++

	if m.fbc != nil {
		m.fbc.OnNavigate()
	}
}

func (m *MenuScene) inputPrev() {
	if m.selection == 0 {
		return
	}

	m.selection--

	if m.fbc != nil {
		m.fbc.OnNavigate()
	}
}

func (m *MenuScene) inputSelect() {
	if m.fbc != nil {
		m.fbc.OnSelect()
	}
	m.choiceMade = true
}
