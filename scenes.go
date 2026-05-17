package main

import (
	"github.com/Zyko0/go-sdl3/sdl"
)

type (
	UserIntent uint8
	SceneEvent uint8
)

const (
	IntentNext UserIntent = iota
	IntentPrev
	IntentUp
	IntentDown
	IntentSelect
	IntentBack
)

const (
	EventWallClockMinutePassed SceneEvent = iota
	EventNetworkChanged
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
	// After calling Input, a redraw is automatically issued
	Input(intent UserIntent)

	// Event is called for other events than user intention
	// After calling Event, a redraw is automatically issued
	// Additional info can be passed through data
	Event(event SceneEvent, data any)

	// MaxInhibitMS is the maximum time this scene allows inhibition,
	// accessed just before inhibition,
	// positive: maximum inhibition time
	// zero = no inhibition allowed, even if the Draw fuction returns false
	// negative: inhibit forever
	MaxInhibitMS() int32
}
