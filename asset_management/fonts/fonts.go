package fonts

import (
	"fmt"
	"log"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
)

type PreRenderedTextKey struct {
	FontID int
	Text   string
}

var (
	fonts = make(map[int]*ttf.Font)

	texts = make(map[PreRenderedTextKey]*ttf.Text) // TODO: currently this WILL grow indefinitely!!!! (maybe not an issue in such a short lived application for now)

	textEngine *ttf.TextEngine
)

func LoadFont(id int, path string, ptSize float32) error {
	font, err := ttf.OpenFont(path, ptSize)
	if err != nil {
		log.Println("Failed to load font")
		return err
	}

	fonts[id] = font

	log.Printf("%s -> %d loaded, family: %+v", path, id, font.FamilyName())
	return nil
}

func Init(renderer *sdl.Renderer) error {
	var err error
	textEngine, err = ttf.CreateRendererTextEngine(renderer)
	return err
}

func GetFont(id int) (*ttf.Font, bool) {
	surface, ok := fonts[id]
	return surface, ok // maybe crash?
}

func GetText(fontID int, text string) (*ttf.Text, error) {
	font, ok := GetFont(fontID)
	if !ok {
		return nil, fmt.Errorf("invalid font id: %d", fontID)
	}

	key := PreRenderedTextKey{
		FontID: fontID,
		Text:   text,
	}

	preRenderedText, ok := texts[key]
	if ok {
		return preRenderedText, nil
	}

	renderedText, err := textEngine.CreateText(font, text)
	if err != nil {
		return nil, err
	}

	texts[key] = renderedText
	return renderedText, nil
}

func Close() {
	for _, preRenderedText := range texts {
		preRenderedText.Destroy()
	}
	texts = nil // ensure the GC can clean it up

	for _, font := range fonts {
		font.Close()
	}
	fonts = nil

	textEngine.DestroyRenderer()
	textEngine = nil
}
