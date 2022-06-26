package equipment

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/library"
)

const (
	cellSize float64 = 16
)

var (
	equipments []*Equipment
	background *aseprite.Cell
)

func loadBackground() {
	layer, err := library.Layer("frame", "base")
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
	xOffset := float64(0)
	yOffset := float64(0)
	if background == nil {
		loadBackground()
	}
	for i := 0; i < 16; i++ {
		xOffset = float64(i) * float64(cellSize)
		if xOffset >= cellSize*8 {
			yOffset = 32
			xOffset = 0
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(background.PositionX), float64(background.PositionY))
		op.GeoM.Translate(xOffset, yOffset+16)
		//op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())
		screen.DrawImage(background.EbitenImage, op)

		if len(equipments) > i {
			err = equipments[i].Draw(xOffset, yOffset+16, screen)
			if err != nil {
				return fmt.Errorf("draw: %w", err)
			}
		}

	}
	return nil
}

func Clear() {
	equipments = []*Equipment{}
}

// Set will set the level of a specific weapon type
func Set(weaponType int, level int) {
	for _, e := range equipments {
		if e.data.Bullet.SourceWeaponType != weaponType {
			continue
		}
		e.level = level
		return
	}
}
