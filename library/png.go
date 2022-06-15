package library

import (
	"fmt"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func loadPng(assetName string) (*ebiten.Image, error) {

	r, err := ReadFile(assetName)
	if err != nil {
		return nil, fmt.Errorf("readFile: %w", err)
	}

	img, err := png.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("png decode: %w", err)
	}

	return ebiten.NewImageFromImage(img), nil
}
