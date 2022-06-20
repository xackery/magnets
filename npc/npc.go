package npc

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/collision"
	"github.com/xackery/magnets/entity"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

type Npc struct {
	hid        string
	entityID   uint
	layer      *aseprite.Layer
	image      *ebiten.Image
	maxHP      int
	hp         int
	spriteName string
	layerName  string
	x          float64
	y          float64
	animation  animation
	moveSpeed  float64
}

type animation struct {
	tag              *aseprite.Tag
	delay            time.Time
	index            int16
	isPingPongToggle bool
}

func New(npcType int, x float64, y float64) (*Npc, error) {
	npcData, ok := npcTypes[npcType]
	if !ok {
		return nil, fmt.Errorf("npc type %d not found", npcType)
	}

	name := "base"
	layer, err := library.Layer(npcData.Sprite.spriteName, npcData.Sprite.layerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", name)
	}

	n := &Npc{
		spriteName: npcData.Sprite.spriteName,
		layerName:  npcData.Sprite.layerName,
		hp:         npcData.MaxHP,
		maxHP:      npcData.MaxHP,
		moveSpeed:  npcData.MoveSpeed,
		x:          x,
		y:          y,
		layer:      layer,
		entityID:   entity.NextEntityID(),
		image:      layer.Cells[0].EbitenImage,
	}

	err = n.SetAnimation("walk")
	if err != nil {
		return nil, fmt.Errorf("SetAnimation %s: %w", "walk", err)
	}

	npcs = append(npcs, n)
	err = entity.Register(n)
	if err != nil {
		return nil, fmt.Errorf("entity.Register: %w", err)
	}
	return n, nil
}

func (n *Npc) IsHit(x, y float64) bool {
	if n.IsDead() {
		return false
	}
	return n.image.At(int(x-n.x), int(y-n.y)).(color.RGBA).A > 0
}

func (n *Npc) Draw(screen *ebiten.Image) error {
	if n.IsDead() {
		return nil
	}
	n.animationStep()
	if len(n.layer.Cells) <= int(n.animation.index) {
		return fmt.Errorf("animationIndex %d is out of bounds for body cells %d", n.animation.index, len(n.layer.Cells))
	}
	c := n.layer.Cells[int(n.animation.index)]

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(n.layer.SpriteWidth/2), -float64(n.layer.SpriteHeight/2))
	op.GeoM.Translate(camera.X, camera.Y)
	op.GeoM.Translate(float64(c.PositionX), float64(c.PositionY))
	op.GeoM.Translate(n.x, n.y)
	op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())

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
	/*x := n.x
	y := n.y + 30
	x -= float64(n.SWidth() / 2)
	y += float64(n.SHeight() / 2)*/

	return nil
}

func (n *Npc) Update(playerX, playerY float64) {

	if !isAIEnabled {
		return
	}

	x := playerX - n.x
	y := playerY - n.y

	if x > 0 && x > n.moveSpeed {
		x = n.moveSpeed
	}
	if x < 0 && -x > -n.moveSpeed {
		x = -n.moveSpeed
	}

	if y > 0 && y > n.moveSpeed {
		y = n.moveSpeed
	}
	if y < 0 && -y > -n.moveSpeed {
		y = -n.moveSpeed
	}

	targetX := n.x + x
	targetY := n.y + y
	if collision.Image.At(int(targetX), int(targetY)).(color.RGBA).A == 0 {
		n.x += x
		n.y += y
		return
	}

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

func (n *Npc) SetPosition(x, y float64) {
	n.x, n.y = x, y
}

func (n *Npc) Position() (float64, float64) {
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

func (n *Npc) X() float64 {
	return n.x
}

func (n *Npc) Y() float64 {
	return n.y
}
