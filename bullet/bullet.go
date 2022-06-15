package bullet

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

type Bullet struct {
	hid        string
	entityID   uint
	layer      *aseprite.Layer
	image      *ebiten.Image
	damage     int
	anchor     int
	spriteName string
	layerName  string
	x          float32
	y          float32
	key        int
	xOffset    float32
	yOffset    float32
	animation  animation
	isDead     bool
}

type animation struct {
	tag              *aseprite.Tag
	delay            time.Time
	index            int16
	isPingPongToggle bool
}

func New(bulletType int, key int, anchor int, xOffset float32, yOffset float32) (*Bullet, error) {

	bType, ok := bulletTypes[bulletType]
	if !ok {
		return nil, fmt.Errorf("bullet of type %d not found", bulletType)
	}
	name := "base"
	layer, err := library.Layer(bType.spriteName, bType.layerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", name)
	}

	n := &Bullet{
		spriteName: bType.spriteName,
		layerName:  bType.layerName,
		damage:     bType.damage,
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

	bullets = append(bullets, n)
	err = entity.Register(n)
	if err != nil {
		return nil, fmt.Errorf("entity.Register: %w", err)
	}
	return n, nil
}

func (n *Bullet) IsHit(x, y float32) bool {
	if n.IsDead() {
		return false
	}
	if n.x > float32(x) {
		return false
	}
	if n.x+float32(n.layer.SpriteWidth) < float32(x) {
		return false
	}
	if n.y > float32(y) {
		return false
	}
	if n.y+float32(n.layer.SpriteHeight) < float32(y) {
		return false
	}
	return true
	//return s.image.At(x-s.x, y-s.y).(color.RGBA).A > 0
}

func (n *Bullet) SetOffset(x, y float32) {
	n.xOffset = x
	n.yOffset = y
	n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
}

func (n *Bullet) Draw(screen *ebiten.Image) error {
	n.animationStep()
	if len(n.layer.Cells) <= int(n.animation.index) {
		return fmt.Errorf("animationIndex %d is out of bounds for body cells %d", n.animation.index, len(n.layer.Cells))
	}
	c := n.layer.Cells[int(n.animation.index)]

	n.x += 0.5

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(n.layer.SpriteWidth/2), -float64(n.layer.SpriteHeight/2))
	op.GeoM.Translate(float64(c.PositionX), float64(c.PositionY))
	op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())
	op.GeoM.Translate(float64(n.x), float64(n.y))

	if n.IsDead() {
		op.ColorM.Scale(1, 1, 1, 0.6)
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
	x := n.x
	y := n.y + 30
	x -= float32(n.SWidth() / 2)
	y += float32(n.SHeight() / 2)

	return nil
}

// SetAnimation sets the animation of the bullet
func (n *Bullet) SetAnimation(name string) error {
	name = strings.ToLower(name)
	tag, err := library.Tag(n.spriteName, name)
	if err != nil {
		return fmt.Errorf("tag: %w", err)
	}
	n.animation.tag = tag
	n.animation.index = tag.From
	return nil
}

func (n *Bullet) animationStep() {
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

func (n *Bullet) SetPosition(x, y float32) {
	n.x, n.y = x, y
}

func (n *Bullet) Position() (float32, float32) {
	return n.x, n.y
}

func (n *Bullet) HID() string {
	return n.hid
}

func (n *Bullet) IsDead() bool {
	return n.isDead
}

func (n *Bullet) EntityID() uint {
	return n.entityID
}

func (n *Bullet) SWidth() int {
	return int(n.layer.SpriteWidth * uint16(global.ScreenScaleX()))
}

func (n *Bullet) SHeight() int {
	return int(n.layer.SpriteHeight * uint16(global.ScreenScaleY()))
}

func (n *Bullet) AnimationAttack() error {
	return nil
}

func (n *Bullet) AnimationGotHit() error {
	return nil
}
