package image

import (
	"fmt"
	"log"

	"github.com/Zyko0/go-sdl3/img"
	"github.com/Zyko0/go-sdl3/sdl"
)

var (
	textures = make(map[int]*sdl.Texture)

	renderer *sdl.Renderer // needed for converting surface to textures
)

func LoadImageResource(id int, path string) error {
	_, ok := textures[id]
	if ok {
		return fmt.Errorf("image id %d already in use", id)
	}

	surface, err := img.Load(path)
	if err != nil {
		log.Println("ERROR: Failed to load image:", err)
		return err
	}

	if renderer == nil {
		panic("forgot to init")
	}

	textures[id], err = renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Println("ERROR: failed to convert surface to texture: ", err)
		return err
	}

	surface.Destroy()

	log.Printf("%s -> %d loaded, fmt: %+v", path, id, textures[id])
	return nil
}

func Init(r *sdl.Renderer) {
	renderer = r
}

func GetTexture(id int) (*sdl.Texture, bool) {
	surface, ok := textures[id]
	return surface, ok // maybe crash?
}

func Close() {
	for _, texture := range textures {
		texture.Destroy()
	}
	renderer = nil
}
