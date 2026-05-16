package main

import (
	"github.com/Zyko0/go-sdl3/sdl"
)

type UserIntent uint8

const (
	IntentNext   = iota
	IntentPrev   = iota
	IntentUp     = iota
	IntentDown   = iota
	IntentSelect = iota
	IntentBack   = iota
)

// SceneController facilitates chaining scenes, if the next scene is registered and the current scene stops requesting drawing then it starts the next one
type SceneController interface {
	RegisterNextScene(Scene)
}

// FeedbackController plays audio basically
type FeedbackController interface {
	OnNavigate()
	OnSelect()
	OnError()
}

type Scene interface {
	// Draw should return true if the next frame should be drawn as well
	// firstFrame is true when drawing the very first frame
	// dtNS is the delta in nanoseconds, i.e.: the time it took to render the last frame, sometimes it is faked
	// durationNS is the time since this scene is being shown
	// Return:
	// - keepDrawing bool: if this was the last frame before inhibition is needed, rendering intended to be inhibited until the next event if set to false
	// - error
	Draw(renderer *sdl.Renderer, scrW, scrH int, firstFrame bool, dtNS, durationNS uint64) (bool, error)

	// Bind is called when the Scene is "loaded"
	Bind(FeedbackController, SceneController)

	// UnBind is called when the Scene is being "unloaded"
	UnBind()

	// Input handles the user input translated from the different input devices
	Input(intent UserIntent)

	// MaxInhibitMS is the maximum time this scene allows inhibition,
	// accessed just before inhibition,
	// positive: maximum inhibition time
	// zero = no inhibition allowed, even if the Draw fuction returns false
	// negative: inhibit forever
	MaxInhibitMS() int32
}
