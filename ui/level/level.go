package level

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/font"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
	"github.com/xackery/magnets/weapon"
)

type Level struct {
	Data      *weapon.WeaponData
	layer     *aseprite.Layer
	cell      *aseprite.Cell
	x         float64
	y         float64
	width     float64
	height    float64
	isVisible bool
	isDead    bool
	percent   float64
}

func New(w *weapon.WeaponData) (*Level, error) {
	if w == nil {
		return nil, fmt.Errorf("weapon instance not found: %v", w)
	}
	if w.Icon == nil {
		return nil, fmt.Errorf("icon not found on weapon %v", w)
	}
	layer, err := library.Layer(w.Icon.SpriteName, w.Icon.LayerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", w.Icon.LayerName)
	}

	b := &Level{
		layer:  layer,
		width:  float64(layer.SpriteWidth),
		height: float64(layer.SpriteHeight),
		Data:   w,
		//percent:   0,
		isVisible: true,
	}

	cell := layer.Cells[0]
	b.cell = cell
	//b.swidth = b.width * int(global.ScreenScaleX())
	//b.sheight = b.height * int(global.ScreenScaleY())
	//b.xOffset -= b.width / 2
	levels = append(levels, b)
	return b, nil
}

func (n *Level) Draw(screen *ebiten.Image, yOffset float64) error {
	if !n.isVisible {
		return nil
	}

	n.x = float64(n.cell.PositionX) + float64(global.ScreenWidth())/2 - 100
	n.y = float64(n.cell.PositionY) + 100 + float64(yOffset)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(n.x, n.y)
	screen.DrawImage(n.cell.EbitenImage, op)

	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(n.x, n.y)
	screen.DrawImage(n.cell.EbitenImage, op)

	//	op.ColorM.Reset()
	//	op.GeoM.Reset()
	//op.GeoM.Scale(n.width*n.percent, n.height)
	//op.GeoM.Translate(n.x, n.y)
	//op.GeoM.Translate(float64(n.cell.PositionX), float64(n.cell.PositionY))
	//	screen.DrawImage(n.cell.EbitenImage, op)

	//bounds := text.BoundString(font.TinyFont(), n.data.Name())
	//x, y := int(n.width/2)-bounds.Min.X-bounds.Dx()/2, int(n.height/2)-bounds.Min.Y-bounds.Dy()/2

	x := n.x + 50
	y := n.y + 10
	text.Draw(screen, n.Data.Name(), font.TinyFont(), int(x-1), int(y-1), color.Black)
	text.Draw(screen, n.Data.Name(), font.TinyFont(), int(x), int(y), color.White)
	return nil
}

func (n *Level) SetDead(state bool) {
	n.isDead = state
}

func (n *Level) SetPosition(x, y float64) {
	n.x = x
	n.y = y
}

func (n *Level) SetPercent(value float64) {
	if value > 1 {
		value = 1
	}
	if value < 0 {
		value = 0
	}
	n.percent = value
}

func (n *Level) SetWidth(value float64) {
	if value == 0 {
		n.width = float64(n.layer.SpriteWidth)
		return
	}
	n.width = value
}

func (n *Level) SetHeight(value float64) {
	if value == 0 {
		n.height = float64(n.layer.SpriteHeight)
		return
	}
	n.height = value
}

func (n *Level) SetVisible(value bool) {
	n.isVisible = value
}

func (n *Level) IsHit(x, y float64) bool {
	if !n.isVisible {
		return false
	}

	if n.x > x {
		return false
	}
	if n.x+float64(n.layer.SpriteWidth)+200 < x {
		return false
	}
	if n.y > y {
		return false
	}
	if n.y+float64(n.layer.SpriteHeight) < y {
		return false
	}
	return true
}
