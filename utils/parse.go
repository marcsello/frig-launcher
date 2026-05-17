package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Zyko0/go-sdl3/sdl"
)

func ParseResolution(res string) (int32, int32, error) {
	parts := strings.SplitN(res, "x", 2)

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format")
	}

	w, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	h, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	if w <= 0 || h <= 0 {
		return 0, 0, fmt.Errorf("zero or negative resolution")
	}

	return int32(w), int32(h), nil
}

func ParseVSync(vsync string) (int32, error) {
	// https://wiki.libsdl.org/SDL3/SDL_SetRenderVSync

	switch vsync {
	case "adaptive":
		return sdl.WINDOW_SURFACE_VSYNC_ADAPTIVE, nil
	case "off":
		return sdl.WINDOW_SURFACE_VSYNC_DISABLED, nil
	default:
		value, err := strconv.Atoi(vsync)
		if err != nil {
			return 0, err
		}
		if value < 1 {
			return 0, fmt.Errorf("invalid value for vsync")
		}
		return int32(value), nil
	}
}
