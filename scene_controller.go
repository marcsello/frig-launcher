package main

import (
	"errors"
	"log"

	"github.com/Zyko0/go-sdl3/sdl"
)

type SceneControllerImpl struct {
	currentScene Scene
	nextScene    Scene

	DesiredFPS uint64

	exitRequested bool

	gamepads []*sdl.Gamepad
}

func (s *SceneControllerImpl) desiredDeltaNS() uint64 {
	return 1_000_000_000 / s.DesiredFPS //desired time b/w frames
}

func (s *SceneControllerImpl) Run(renderer *sdl.Renderer, feedbackController FeedbackController, scrW, scrH int) error {
	// Cleanup
	defer func() {
		for _, gamepad := range s.gamepads {
			gamepad.Close()
		}
	}()

	// Main loop
	delta := s.desiredDeltaNS()

	firstFrame := true
	var currentSceneStart uint64

	return sdl.RunLoop(func() error {
		var err error
		loopStarted := sdl.TicksNS()

		if s.exitRequested {
			return sdl.EndLoop
		}

		// In case this is the very first loop (or we lost the current scene somehow)
		if s.currentScene == nil && s.nextScene != nil {
			// Do the initial scene loading
			log.Println("Switching to the first scene...")
			s.currentScene = s.nextScene
			s.nextScene = nil
			s.currentScene.Bind(feedbackController, s)
			firstFrame = true
		}

		var keepDrawing bool
		if s.currentScene != nil {
			if firstFrame {
				currentSceneStart = loopStarted
			}
			keepDrawing, err = s.currentScene.Draw(renderer, scrW, scrH, firstFrame, delta, loopStarted-currentSceneStart)
			if err != nil {
				if errors.Is(err, sdl.EndLoop) {
					log.Println("Termination request from the scene")
				} else {
					log.Println("ERROR: drawing scene:", err)
				}
				return err
			}
			firstFrame = false
		} else {
			err = renderer.SetDrawColor(255, 0, 0, 255)
			if err != nil {
				return err
			}
			err = renderer.DebugText(50, 50, "no scene loaded!")
			if err != nil {
				return err
			}
		}

		err = renderer.Present()
		if err != nil {
			log.Println("Error: present:", err)
			return err
		}

		if !keepDrawing && s.nextScene != nil {
			// do scene switch
			log.Println("Switching scene...")
			s.currentScene.UnBind()
			s.currentScene = s.nextScene
			s.nextScene = nil
			s.currentScene.Bind(feedbackController, s)
			firstFrame = true

			keepDrawing = true // don't block
		}

		s.handleInput(!keepDrawing)

		// calculate if maybe we need to pad the frame time
		delta = sdl.TicksNS() - loopStarted
		desiredDeltaNS := s.desiredDeltaNS()
		if delta < desiredDeltaNS {
			sdl.DelayNS(desiredDeltaNS - delta)
		}

		// update delta with the valid delta time
		delta = sdl.TicksNS() - loopStarted

		// if drawing is inhibited, we may have to fake the delta time
		if !keepDrawing && delta > desiredDeltaNS*2 {
			delta = desiredDeltaNS
		}

		return nil
	})
}

func (s *SceneControllerImpl) proxyInput(intent UserIntent) {
	if s.currentScene != nil {
		s.currentScene.Input(intent)
	} else {
		log.Println("Intent dropped: ", intent)
	}
}

func (s *SceneControllerImpl) handleInput(blocking bool) {
	// This is a spaghetti lol
	var event sdl.Event

	maxInhibitMS := s.currentScene.MaxInhibitMS()

	if maxInhibitMS == 0 {
		// disallow blocking if max inhibit is 0
		blocking = false
	}

	// If we are doing a blocking wait, then wait for the *first* event indefinitely
	if blocking {
		if maxInhibitMS < 0 {
			var err error
			err = sdl.WaitEvent(&event)
			if err != nil {
				log.Println("Error while waiting for event:", err)
				return
			}
		} else {
			gotEvent := sdl.WaitEventTimeout(&event, maxInhibitMS)
			if !gotEvent {
				return
			}
		}
	}

	firstLoop := true
	for {
		if !(firstLoop && blocking) {
			if !sdl.PollEvent(&event) {
				break
			}
		}
		firstLoop = false

		if event.Type == sdl.EVENT_QUIT {
			s.exitRequested = true
			return // any further event processing is kinda useless
		}

		// keyboard
		if event.Type == sdl.EVENT_KEY_DOWN {
			switch event.KeyboardEvent().Key {
			case sdl.K_ESCAPE:
				s.exitRequested = true
				return // any further event processing is kinda useless
			case sdl.K_UP:
				s.proxyInput(IntentUp)
			case sdl.K_DOWN:
				s.proxyInput(IntentDown)
			case sdl.K_LEFT:
				s.proxyInput(IntentPrev)
			case sdl.K_RIGHT:
				s.proxyInput(IntentNext)
			case sdl.K_RETURN, sdl.K_RETURN2, sdl.K_SPACE:
				s.proxyInput(IntentSelect)
			}
		}

		if event.Type == sdl.EVENT_GAMEPAD_BUTTON_DOWN {
			switch sdl.GamepadButton(event.GamepadButtonEvent().Button) {
			case sdl.GAMEPAD_BUTTON_DPAD_UP:
				s.proxyInput(IntentUp)
			case sdl.GAMEPAD_BUTTON_DPAD_DOWN:
				s.proxyInput(IntentDown)
			case sdl.GAMEPAD_BUTTON_DPAD_LEFT:
				s.proxyInput(IntentPrev)
			case sdl.GAMEPAD_BUTTON_DPAD_RIGHT:
				s.proxyInput(IntentNext)
			case sdl.GAMEPAD_BUTTON_BACK, sdl.GAMEPAD_BUTTON_EAST:
				s.proxyInput(IntentBack)
			case sdl.GAMEPAD_BUTTON_START, sdl.GAMEPAD_BUTTON_SOUTH:
				s.proxyInput(IntentSelect)
			}
		}

		if event.Type == sdl.EVENT_GAMEPAD_ADDED {
			joystickID := event.GamepadDeviceEvent().Which
			gamepad, err := joystickID.OpenGamepad()
			if err != nil {
				log.Println("ERROR: new gamepad detected, but could not open it:", err)
			}

			s.gamepads = append(s.gamepads, gamepad)
			log.Println("new gamepad added: ", gamepad.Name())
		}
		if event.Type == sdl.EVENT_GAMEPAD_REMAPPED {
			log.Println("gamepad remapped")
		}
		if event.Type == sdl.EVENT_GAMEPAD_REMOVED {
			log.Println("gamepad removed")
		}
	}
}

func (s *SceneControllerImpl) RegisterNextScene(scene Scene) {
	s.nextScene = scene
}

func (s *SceneControllerImpl) RequestExit() {
	ev := sdl.Event{Type: sdl.EVENT_QUIT}
	err := sdl.PushEvent(&ev) // the docs state: It is safe to call this function from any thread.
	if err != nil {
		log.Println("failed to push quit event")
	}
}
