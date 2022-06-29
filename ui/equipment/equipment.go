package equipment

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/library"
	"github.com/xackery/magnets/weapon"
)

type Equipment struct {
	layer     *aseprite.Layer
	cell      *aseprite.Cell
	data      *weapon.WeaponData
	isVisible bool
	isDead    bool
	level     int
	text      string
}

func New(data *weapon.WeaponData) (*Equipment, error) {
	if len(equipments) >= 16 {
		return nil, fmt.Errorf("already full of equipment")
	}
	layer, err := library.Layer(data.Icon.SpriteName, data.Icon.LayerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", data.Icon.LayerName)
	}

	b := &Equipment{
		layer: layer,
		//percent:   0,
		isVisible: true,
	}

	cell := layer.Cells[0]
	b.cell = cell
	//b.swidth = b.width * int(global.ScreenScaleX())
	//b.sheight = b.height * int(global.ScreenScaleY())
	//b.xOffset -= b.width / 2

	equipments = append(equipments, b)
	return b, nil
}

func (n *Equipment) Draw(xOffset float64, yOffset float64, screen *ebiten.Image) error {
	if !n.isVisible {
		return nil
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(n.cell.PositionX), float64(n.cell.PositionY))
	//op.GeoM.Scale(n.width, n.height)
	op.GeoM.Translate(xOffset, yOffset)
	//op.ColorM.Scale(0.7, 0.7, 0.7, 0.7)
	//op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())
	screen.DrawImage(n.cell.EbitenImage, op)

	//bounds := text.BoundString(font.NormalFont(), n.text)
	//x, y := int(n.width/2)-bounds.Min.X-bounds.Dx()/2, int(n.height/2)-bounds.Min.Y-bounds.Dy()/2

	//text.Draw(screen, n.text, font.NormalFont(), x-1, y-1, color.Black)
	//text.Draw(screen, n.text, font.NormalFont(), x, y, color.White)
	return nil
}

func (n *Equipment) SetDead(state bool) {
	n.isDead = state
}

func (n *Equipment) SetLevel(value int) {
	n.level = value
}

func (n *Equipment) SetText(value string) {
	n.text = value
}
