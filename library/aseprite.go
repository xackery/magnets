package library

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/aseprite"
)

var (
	sprites = make(map[string]*aseprite.Sprite)
)

func loadAseprite(assetName string) (*aseprite.Sprite, error) {
	var sprite *aseprite.Sprite

	r, err := ReadFile(assetName)
	if err != nil {
		return nil, fmt.Errorf("readFile: %w", err)
	}

	sprite, err = aseprite.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("aseprite decode: %w", err)
	}

	uniqueNames := []string{}
	//convert all images to ebiten format
	for _, l := range sprite.Layers {
		for _, un := range uniqueNames {
			if l.Name == un {
				return nil, fmt.Errorf("%s is a duplicate layer", l.Name)
			}
		}
		for cellIndex, c := range l.Cells {
			assessPivots(assetName, l, c, cellIndex)
			c.EbitenImage = ebiten.NewImageFromImage(c.Image)
			if err != nil {
				return nil, fmt.Errorf("newImageFromEbiten: %w", err)
			}
			c.Image = nil
		}
	}
	return sprite, nil
}

func SpriteConstDump() map[string]*aseprite.Sprite {
	return sprites
}
