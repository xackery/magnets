package world

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/global"
)

var (
	worlds []*World
)

func init() {
	global.SubscribeOnResolutionChange(onResolutionChange)
}

func At(x, y float64) *World {
	for _, p := range worlds {
		if !p.IsHit(x, y) {
			continue
		}
		return p
	}
	return nil
}

func Draw(screen *ebiten.Image) error {
	var err error
	for _, p := range worlds {
		err = p.Draw(screen)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}

	return nil
}

func Update() {
	for _, p := range worlds {
		p.Update()
	}
}

func onResolutionChange() {
	/*for _, n := range worlds {
		n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	}*/

	camera.X = float64(global.ScreenWidth() / 4)
	camera.Y = float64(global.ScreenHeight() / 4)
	if len(worlds) > 0 {
		camera.X -= worlds[0].x
		camera.Y -= worlds[0].y
	}

}

func Worlds() []*World {
	return worlds
}

func HID(hid string) *World {
	for _, n := range worlds {
		if n.HID() == hid {
			return n
		}
	}
	return nil
}

func Clear() {
	worlds = []*World{}
}

func X() float64 {
	if len(worlds) == 0 {
		return 0
	}
	return worlds[0].x
}

func Y() float64 {
	if len(worlds) == 0 {
		return 0
	}
	return worlds[0].y
}

func IsCollision(x, y float64) bool {
	if len(worlds) == 0 {
		return false
	}
	return worlds[0].IsHit(x, y)
}
