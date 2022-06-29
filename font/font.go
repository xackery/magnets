package font

import (
	"fmt"
	"sync"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	mplusTinyFont   font.Face
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
	mutex           sync.RWMutex
)

// Load loads all fonts
func Load() error {
	mutex.Lock()
	defer mutex.Unlock()
	dpi := float64(72)
	/*
		buf := &bytes.Buffer{}

			r, err := library.ReadFile("cnc.ttf")
			if err != nil {
				return fmt.Errorf("library.ReadFile: %w", err)
			}

			_, err = io.Copy(buf, r)
			if err != nil {
				return fmt.Errorf("copy: %w", err)
			}

			tt, err := truetype.Parse(buf.Bytes())
	*/
	tt, err := truetype.Parse(goregular.TTF)

	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	mplusTinyFont = truetype.NewFace(tt, &truetype.Options{
		Size:    8,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	mplusSmallFont = truetype.NewFace(tt, &truetype.Options{
		Size:    10,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	mplusNormalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	mplusBigFont = truetype.NewFace(tt, &truetype.Options{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	return nil
}

// BigFont returns the big font
func BigFont() font.Face {
	mutex.Lock()
	defer mutex.Unlock()
	return mplusBigFont
}

// NormalFont returns the normal font
func NormalFont() font.Face {
	mutex.Lock()
	defer mutex.Unlock()
	return mplusNormalFont
}

// SmallFont returns the small font
func SmallFont() font.Face {
	mutex.Lock()
	defer mutex.Unlock()
	return mplusSmallFont
}

// TinyFont returns the tiny font
func TinyFont() font.Face {
	mutex.Lock()
	defer mutex.Unlock()
	return mplusTinyFont
}
