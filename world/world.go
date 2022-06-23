package world

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/entity"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

type World struct {
	hid      string
	entityID uint
	image    *ebiten.Image
	layer    *aseprite.Layer
	x        float64
	y        float64
}

func New(worldType int) (*World, error) {
	wData, ok := worldTypes[worldType]
	if !ok {
		return nil, fmt.Errorf("world %d not found", worldType)
	}

	layer, err := library.Layer(wData.Tilemap.spriteName, wData.Tilemap.layerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on sprite %s layer %s", wData.Tilemap.spriteName, wData.Tilemap.layerName)
	}

	n := &World{
		image: layer.Cells[0].EbitenImage,
		layer: layer,
	}
	n.image = ebiten.NewImage(int(layer.SpriteWidth), int(layer.SpriteHeight))

	worlds = append(worlds, n)
	err = entity.Register(n)
	if err != nil {
		return nil, fmt.Errorf("entity.Register: %w", err)
	}
	return n, nil
}

func (n *World) IsHit(x, y float64) bool {

	/*if n.x > float64(x) {
		return false
	}
	if n.x+float64(n.layer.SpriteWidth) < float64(x) {
		return false
	}
	if n.y > float64(y) {
		return false
	}
	if n.y+float64(n.layer.SpriteHeight) < float64(y) {
		return false
	}
	return true*/
	return n.image.At(int(x), int(y)).(color.RGBA).A > 0
	//return s.image.At(x-s.x, y-s.y).(color.RGBA).A > 0
}

func (n *World) Draw(screen *ebiten.Image) error {

	c := n.layer.Cells[0]
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(n.layer.SpriteWidth/2), -float64(n.layer.SpriteHeight/2))
	op.GeoM.Translate(float64(c.PositionX), float64(c.PositionY))
	op.GeoM.Translate(float64(camera.X), float64(camera.Y))
	op.GeoM.Translate(float64(n.x), float64(n.y))
	op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())

	//op.GeoM.Translate(float64(global.ScreenWidth())/2, float64(global.ScreenHeight())/2)
	//op.GeoM.Translate(float64(n.x), float64(n.y))

	//duplication spell effect?
	/*for j := -128; j <= 128; j += 32 {
		//for i := -1; i <= 1; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())
		op.GeoM.Translate(float64(int(c.PositionX)+n.x+n.dragX), float64(int(c.PositionY)+n.y+n.dragY+j))

		//op.GeoM.Translate(1, float64(j))
		// Alpha scale should be 1.0/49.0, but accumulating 1/49 49 times doesn't reach to 1, because
		// the final color is affected by the destination alpha when CompositeModeSourceOver is used.
		// This composite mode is the default mode. See how this is calculated at the doc:
		// https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#CompositeMode
		//
		// Use a higher value than 1.0/49.0. Here, 1.0/25.0 here to get a reasonable result.
		op.ColorM.Scale(1, 1, 1, 0.8)
		screen.DrawImage(c.EbitenImage, op)
		//}
	}*/

	screen.DrawImage(c.EbitenImage, op)
	//text.Draw(screen, n.nameTag, font.TinyFont(), n.x-(len(n.nameTag)*2)+1, n.y+int(n.layer.SpriteHeight)+40+1, color.Black)
	//text.Draw(screen, n.nameTag, font.TinyFont(), n.x-(len(n.nameTag)*2), n.y+int(n.layer.SpriteHeight)+40, color.White)
	//x := n.x
	//y := n.y + 30
	//x -= float64(n.SWidth() / 2)
	//y += float64(n.SHeight() / 2)

	return nil
}

func (n *World) Update() {

}

func (n *World) SetPosition(x, y float64) {
	n.x, n.y = x, y
}

func (n *World) Position() (float64, float64) {
	return n.x, n.y
}

func (n *World) HID() string {
	return n.hid
}

func (n *World) EntityID() uint {
	return n.entityID
}

func (n *World) SWidth() int {
	return int(n.layer.SpriteWidth * uint16(global.ScreenScaleX()))
}

func (n *World) SHeight() int {
	return int(n.layer.SpriteHeight * uint16(global.ScreenScaleY()))
}

func (n *World) AnimationAttack() error {
	return nil
}

func (n *World) AnimationGotHit() error {
	return nil
}

func (n *World) X() float64 {
	return n.x
}

func (n *World) Y() float64 {
	return n.y
}
