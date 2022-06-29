package life

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

const (
	cellSize float64 = 16
	countRow float64 = 10
)

var (
	hp    int
	maxHP int
	base  *aseprite.Cell
	heart *aseprite.Cell
)

func loadHeart() {
	layer, err := library.Layer("heart", "base")
	if err != nil {
		log.Error().Err(err).Msgf("library.layer heart base")
		return
	}
	if len(layer.Cells) < 1 {
		log.Error().Msgf("no cells found on layer %s", "heart")
		return
	}

	base = layer.Cells[0]

	layer, err = library.Layer("heart", "full")
	if err != nil {
		log.Error().Err(err).Msgf("library.layer heart base")
		return
	}
	if len(layer.Cells) < 1 {
		log.Error().Msgf("no cells found on layer %s", "heart")
		return
	}
	heart = layer.Cells[0]

}

func Draw(screen *ebiten.Image) error {
	xOffset := float64(-cellSize)
	yOffset := float64(0)
	if base == nil {
		loadHeart()
	}

	for i := 0; i < 20; i++ {
		xOffset += cellSize
		if xOffset >= cellSize*countRow {
			yOffset = 8
			xOffset = 0
		}

		if maxHP <= i {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(base.PositionX), float64(base.PositionY))
		op.GeoM.Translate((float64(global.ScreenWidth())-cellSize*countRow)+xOffset, yOffset+16)
		screen.DrawImage(base.EbitenImage, op)

		if hp == 0 {
			continue
		}
		if hp <= i {
			continue
		}
		op.GeoM.Reset()
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(float64(heart.PositionX), float64(heart.PositionY))
		op.GeoM.Translate((float64(global.ScreenWidth())-cellSize*countRow)+xOffset, yOffset+16)
		screen.DrawImage(heart.EbitenImage, op)

	}
	return nil
}

func Clear() {
	hp = 0
}

func SetHP(value int) {
	hp = value
}

func SetMaxHP(value int) {
	maxHP = value
}
