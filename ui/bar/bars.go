package bar

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	bars []*Bar
)

func Draw(screen *ebiten.Image) error {
	var err error
	for _, b := range bars {
		err = b.Draw(screen)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}
	return nil
}

func Clear() {
	bars = []*Bar{}
}
