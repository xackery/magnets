package library

import (
	"fmt"

	"github.com/goki/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/language"
)

// FontData represents a font data
type FontData struct {
	Face                font.Face
	Height              int
	Name                string
	Language            language.Tag
	BoundStringCache    map[font.Face]map[string]*BoundStringCacheEntry
	RenderingLineHeight int
}

// BoundStringCacheEntry is used for font boundings
type BoundStringCacheEntry struct {
	bounds  *fixed.Rectangle26_6
	advance fixed.Int26_6
}

// LoadFontTTF instantiates a truetype font. truetype.Parse() can be used to load a TTF to fontData
func LoadFontTTF(name string, fontData []byte, opts *truetype.Options, r rune) error {
	if opts == nil {
		opts = &truetype.Options{Size: 12, DPI: 72, Hinting: font.HintingFull}
	}

	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return errors.Wrap(err, "parse ttf font")
	}
	f := &FontData{
		Name:                name,
		BoundStringCache:    make(map[font.Face]map[string]*BoundStringCacheEntry),
		RenderingLineHeight: 18,
	}
	f.Face = truetype.NewFace(tt, opts)
	b, _, ok := f.Face.GlyphBounds(r)
	if !ok {
		return fmt.Errorf("calibrate glyph bounds failed")
	}
	f.Height = (b.Max.Y - b.Min.Y).Ceil()

	fonts[name] = f
	return nil
}

// Font returns a font represented by name
func Font(name string) *FontData {
	return fonts[name]
}
