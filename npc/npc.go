package npc

import (
	"fmt"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/entity"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

type Npc struct {
	hid        string
	entityID   uint
	layer      *aseprite.Layer
	image      *ebiten.Image
	anchor     int
	maxHP      int
	hp         int
	spriteName string
	layerName  string
	x          float32
	y          float32
	key        int
	xOffset    float32
	yOffset    float32
	animation  animation
	nameTag    string
}

type animation struct {
	tag              *aseprite.Tag
	delay            time.Time
	index            int16
	isPingPongToggle bool
}

func New(spriteName string, layerName string, key int, anchor int, xOffset float32, yOffset float32) (*Npc, error) {

	name := "base"
	layer, err := library.Layer(spriteName, layerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", name)
	}

	n := &Npc{
		spriteName: spriteName,
		layerName:  layerName,
		hp:         1,
		maxHP:      1,
		key:        key,
		anchor:     anchor,
		xOffset:    xOffset,
		yOffset:    yOffset,
		layer:      layer,
		entityID:   entity.NextEntityID(),
		image:      layer.Cells[0].EbitenImage,
	}

	n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	err = n.SetAnimation("left")
	if err != nil {
		return nil, fmt.Errorf("SetAnimation %s: %w", "left", err)
	}

	npcs = append(npcs, n)
	err = entity.Register(n)
	if err != nil {
		return nil, fmt.Errorf("entity.Register: %w", err)
	}
	return n, nil
}

func (n *Npc) IsHit(x, y float32) bool {
	if n.IsDead() {
		return false
	}
	//fmt.Println(x, y, "vs", n.x, n.y, n.layer.SpriteWidth, n.layer.SpriteHeight)
	if n.x-(float32(n.layer.SpriteWidth)/2) > float32(x) {
		return false
	}
	if n.x+(float32(n.layer.SpriteWidth)/2) < float32(x) {
		return false
	}
	if n.y-float32(n.layer.SpriteHeight/2) > float32(y) {
		return false
	}
	if n.y+float32(n.layer.SpriteHeight/2) < float32(y) {
		return false
	}
	return true
	//return s.image.At(x-s.x, y-s.y).(color.RGBA).A > 0
}

func (n *Npc) SetOffset(x, y float32) {
	n.xOffset = x
	n.yOffset = y
	n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
}

func (n *Npc) Draw(screen *ebiten.Image) error {
	n.animationStep()
	if len(n.layer.Cells) <= int(n.animation.index) {
		return fmt.Errorf("animationIndex %d is out of bounds for body cells %d", n.animation.index, len(n.layer.Cells))
	}
	c := n.layer.Cells[int(n.animation.index)]

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(n.layer.SpriteWidth/2), -float64(n.layer.SpriteHeight/2))
	op.GeoM.Translate(float64(c.PositionX), float64(c.PositionY))
	op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())
	op.GeoM.Translate(float64(n.x), float64(n.y))

	if n.IsDead() {
		op.ColorM.Scale(1, 1, 1, 0.6)
	} else if n.hp <= (n.maxHP / 2) {
		op.ColorM.Scale(1, 0.5, 0.5, 1)
	} else {
		op.ColorM.Scale(1, 1, 1, 1)
	}

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

	if n.IsDead() {
		return nil
	}
	//text.Draw(screen, n.nameTag, font.TinyFont(), n.x-(len(n.nameTag)*2)+1, n.y+int(n.layer.SpriteHeight)+40+1, color.Black)
	//text.Draw(screen, n.nameTag, font.TinyFont(), n.x-(len(n.nameTag)*2), n.y+int(n.layer.SpriteHeight)+40, color.White)
	x := n.x
	y := n.y + 30
	x -= float32(n.SWidth() / 2)
	y += float32(n.SHeight() / 2)

	return nil
}

// SetAnimation sets the animation of the npc
func (n *Npc) SetAnimation(name string) error {
	name = strings.ToLower(name)
	tag, err := library.Tag(n.spriteName, name)
	if err != nil {
		return fmt.Errorf("tag: %w", err)
	}
	n.animation.tag = tag
	n.animation.index = tag.From
	return nil
}

func (n *Npc) animationStep() {
	if n.IsDead() {
		return
	}
	if n.animation.delay.After(time.Now()) {
		return
	}

	if n.animation.tag.AnimationDirection == 2 && n.animation.isPingPongToggle {
		n.animation.index--

		if n.animation.index <= n.animation.tag.From {
			n.animation.isPingPongToggle = false
		}
		if n.animation.index < 0 {
			n.animation.index = 0
		}
	} else {
		n.animation.index++
		if n.animation.index > n.animation.tag.To {
			if n.animation.tag.AnimationDirection == 2 {
				n.animation.index -= 2
				n.animation.isPingPongToggle = true
			} else {
				n.animation.index = n.animation.tag.From
			}
		}
	}

	if n.layer == nil {
		return
	}

	if len(n.layer.Cells) <= int(n.animation.index) {
		return
	}

	c := n.layer.Cells[int(n.animation.index)]
	n.animation.delay = time.Now().Add(time.Duration(c.Duration) * time.Millisecond)
}

func (n *Npc) SetPosition(x, y float32) {
	n.x, n.y = x, y
}

func (n *Npc) Position() (float32, float32) {
	return n.x, n.y
}

func (n *Npc) HID() string {
	return n.hid
}

func (n *Npc) Damage(damage int) bool {
	n.hp -= damage
	if n.hp < 1 {
		n.hp = 0
		return true
	}
	return false
}

func (n *Npc) IsDead() bool {
	return n.hp < 1
}

func (n *Npc) EntityID() uint {
	return n.entityID
}

func (n *Npc) SWidth() int {
	return int(n.layer.SpriteWidth * uint16(global.ScreenScaleX()))
}

func (n *Npc) SHeight() int {
	return int(n.layer.SpriteHeight * uint16(global.ScreenScaleY()))
}

func (n *Npc) AnimationAttack() error {
	return nil
}

func (n *Npc) AnimationGotHit() error {
	return nil
}

func (n *Npc) X() float32 {
	return n.x
}

func (n *Npc) Y() float32 {
	return n.y
}
