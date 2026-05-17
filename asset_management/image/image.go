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

// LoadImageResource loads an image resource to a resource id. Parameters svgScaleW and svgScaleH are only used when loading an SVG image, ignored otherwise.
func LoadImageResource(id int, path string, svgScaleW, svgScaleH int32) error {
	_, ok := textures[id]
	if ok {
		return fmt.Errorf("image id %d already in use", id)
	}
	if renderer == nil {
		panic("forgot to init")
	}

	ioStream, err := sdl.IOFromFile(path, "r")
	if err != nil {
		log.Println("Failed to open resource:", err)
		return err
	}
	defer func(ioStream *sdl.IOStream) {
		err := ioStream.Close()
		if err != nil {
			log.Printf("WARNING: failed to close io stream for %s: %v", path, err)
		}
	}(ioStream)

	var surface *sdl.Surface
	if img.IsSVG(ioStream) && (svgScaleH > 0 || svgScaleW > 0) {
		log.Printf("%s is an SVG image, scaling to %dx%d", path, svgScaleW, svgScaleH)
		surface, err = img.LoadSizedSVG_IO(ioStream, svgScaleW, svgScaleH)
	} else {
		surface, err = img.LoadIO(ioStream, false) // we have a deferred close
	}
	if err != nil {
		log.Println("ERROR: Failed to load image:", err)
		return err
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
