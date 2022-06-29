package gameover

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/rs/zerolog/log"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/font"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

const ()

var (
	title      string
	background *aseprite.Cell
	texts      []string
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
	if !global.IsGameOver() {
		return nil
	}
	if background == nil {
		loadBackground()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(background.PositionX), float64(background.PositionY))
	op.GeoM.Translate(float64(global.ScreenWidth())/2-float64(background.EbitenImage.Bounds().Dx()/2), 50)
	if strings.Contains(title, "Game") {
		op.ColorM.Scale(1, 0, 0, 0.7)
	}
	//op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())
	screen.DrawImage(background.EbitenImage, op)

	x := global.ScreenWidth()/2 - 60
	y := 70
	text.Draw(screen, title, font.NormalFont(), x-1, y-1, color.Black)
	text.Draw(screen, title, font.NormalFont(), x, y, color.White)

	y += 32
	for _, t := range texts {
		y += 32
		text.Draw(screen, t, font.NormalFont(), x-1, y-1, color.Black)
		text.Draw(screen, t, font.NormalFont(), x, y, color.White)
	}

	return nil
}

func Clear() {
	texts = []string{}
}

func AddText(value string) {
	texts = append(texts, value)
}
func SetTitle(value string) {
	title = value + "\nPress R to restart"
}
