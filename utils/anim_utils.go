package utils

import (
	"math"

	"github.com/Zyko0/go-sdl3/sdl"
)

func EaseInOutCubic(x float64) float64 {
	if x < 0.5 {
		return 4 * x * x * x
	}

	return 1 - math.Pow(-2*x+2, 3)/2
}

type TransitionAnimator struct {
	DesiredRect sdl.FRect
	Rect        sdl.FRect
	Active      bool
}

func (p *TransitionAnimator) calc(curVal, desiredVal, corrector, factor float32) float32 {
	dist := curVal - desiredVal
	curVal -= (dist / factor) * corrector

	// prevent overshoot, it's very rare, but lagging can cause it
	if math.Abs(float64(curVal-desiredVal)) > math.Abs(float64(dist)) {
		curVal = desiredVal
	}

	return curVal
}

func (p *TransitionAnimator) Update(dtNS uint64) bool {
	if !p.Active {
		return false
	}

	corrector := float32(dtNS) / 100_000_000 // ~0.33 for 30fps and ~0.16 for 60fps

	p.Rect.X = p.calc(p.Rect.X, p.DesiredRect.X, corrector, 4)
	p.Rect.Y = p.calc(p.Rect.Y, p.DesiredRect.Y, corrector, 4)
	p.Rect.W = p.calc(p.Rect.W, p.DesiredRect.W, corrector, 2)
	p.Rect.H = p.calc(p.Rect.H, p.DesiredRect.H, corrector, 2)

	if (math.Abs(float64(p.Rect.X-p.DesiredRect.X)) < 1) &&
		(math.Abs(float64(p.Rect.Y-p.DesiredRect.Y)) < 1) &&
		(math.Abs(float64(p.Rect.W-p.DesiredRect.W)) < 1) &&
		(math.Abs(float64(p.Rect.H-p.DesiredRect.H)) < 1) {
		p.Active = false
	}

	return p.Active
}

func (p *TransitionAnimator) SetDesired(newRect sdl.FRect) {
	p.DesiredRect = newRect
	p.Active = true
}

func (p *TransitionAnimator) SendTo(x, y float32) {
	p.DesiredRect.X = x
	p.DesiredRect.Y = y
	p.Active = true
}
