package bar

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/font"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

type Bar struct {
	layer     *aseprite.Layer
	cell      *aseprite.Cell
	x         float64
	y         float64
	width     float64
	height    float64
	isVisible bool
	isDead    bool
	percent   float64
	text      string
}

func New(spriteName string, layerName string) (*Bar, error) {
	layer, err := library.Layer(spriteName, layerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", layerName)
	}

	b := &Bar{
		layer:  layer,
		width:  float64(layer.SpriteWidth),
		height: float64(layer.SpriteHeight),
		//percent:   0,
		isVisible: true,
	}

	cell := layer.Cells[0]
	b.cell = cell
	//b.swidth = b.width * int(global.ScreenScaleX())
	//b.sheight = b.height * int(global.ScreenScaleY())
	//b.xOffset -= b.width / 2
	bars = append(bars, b)
	return b, nil
}

func (n *Bar) Draw(screen *ebiten.Image) error {
	if !n.isVisible {
		return nil
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(n.cell.PositionX), float64(n.cell.PositionY))
	op.GeoM.Scale(n.width, n.height)
	op.GeoM.Translate(n.x, n.y)
	op.ColorM.Scale(0.7, 0.7, 0.7, 0.7)
	screen.DrawImage(n.cell.EbitenImage, op)

	op.ColorM.Reset()
	op.GeoM.Reset()
	op.GeoM.Scale(n.width*n.percent, n.height)
	//op.GeoM.Translate(n.x, n.y)
	//op.GeoM.Translate(float64(n.cell.PositionX), float64(n.cell.PositionY))
	screen.DrawImage(n.cell.EbitenImage, op)

	bounds := text.BoundString(font.TinyFont(), n.text)
	x, y := int(n.width/2)-bounds.Min.X-bounds.Dx()/2, int(n.height/2)-bounds.Min.Y-bounds.Dy()/2

	text.Draw(screen, n.text, font.TinyFont(), x-1, y-1, color.Black)
	text.Draw(screen, n.text, font.TinyFont(), x, y, color.White)

	text.Draw(screen, fmt.Sprintf("%d kills", global.Kill), font.TinyFont(), 50+x-1, y+22-1, color.Black)
	text.Draw(screen, fmt.Sprintf("%d kills", global.Kill), font.TinyFont(), 50+x, y+22, color.White)
	return nil
}

func (n *Bar) SetDead(state bool) {
	n.isDead = state
}

func (n *Bar) SetPosition(x, y float64) {
	n.x = x
	n.y = y
}

func (n *Bar) SetPercent(value float64) {
	if value > 1 {
		value = 1
	}
	if value < 0 {
		value = 0
	}
	n.percent = value
}

func (n *Bar) SetWidth(value float64) {
	if value == 0 {
		n.width = float64(n.layer.SpriteWidth)
		return
	}
	n.width = value
}

func (n *Bar) SetHeight(value float64) {
	if value == 0 {
		n.height = float64(n.layer.SpriteHeight)
		return
	}
	n.height = value
}

func (n *Bar) SetText(value string) {
	n.text = value
}
