package level

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/rs/zerolog/log"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/font"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

const (
	title = "Choose an upgrade"
)

var (
	levels     []*Level
	background *aseprite.Cell
)

func loadBackground() {
	layer, err := library.Layer("levelup", "base")
	if err != nil {
		log.Error().Err(err).Msgf("library.layer frame base")
		return
	}
	if len(layer.Cells) < 1 {
		log.Error().Msgf("no cells found on layer %s", "frame")
		return
	}

	background = layer.Cells[0]
}

func Draw(screen *ebiten.Image) error {
	var err error
	if background == nil {
		loadBackground()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(background.PositionX), float64(background.PositionY))
	op.GeoM.Translate(float64(global.ScreenWidth())/2-float64(background.EbitenImage.Bounds().Dx()/2), 50)
	//op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())
	screen.DrawImage(background.EbitenImage, op)

	text.Draw(screen, title, font.TinyFont(), global.ScreenWidth()/2-1-60, 70-1, color.Black)
	text.Draw(screen, title, font.TinyFont(), global.ScreenWidth()/2-60, 70, color.White)

	for i, b := range levels {
		err = b.Draw(screen, float64(i)*64)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}
	return nil
}

func Clear() {
	levels = []*Level{}
}

func IsHit(x, y float64) *Level {
	for _, n := range levels {
		if n.IsHit(x, y) {
			return n
		}
	}
	return nil
}
func ByIndex(index int) *Level {
	if len(levels) <= index {
		return nil
	}
	return levels[index]
}
