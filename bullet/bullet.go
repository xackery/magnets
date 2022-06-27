package bullet

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/entity"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

type Bullet struct {
	hid          string
	entityID     uint
	layer        *aseprite.Layer
	image        *ebiten.Image
	behaviorType int
	direction    int
	lifespan     time.Time
	data         *BulletData
	spriteName   string
	layerName    string
	rotation     float64
	x            float64
	y            float64
	spawnX       float64
	spawnY       float64
	player       entity.Entiter
	distance     float64 //distance an object should travel to
	animation    animation
	isDead       bool
	isImmortal   bool //will the bullet die on impact
	isReturning  bool // For behaviors that go to and fro, this is the returning flag
	moveSpeed    float64
	spawnCount   int // number of bullets of type spawned when this initially appeared
}

type animation struct {
	tag              *aseprite.Tag
	delay            time.Time
	index            int16
	isPingPongToggle bool
}

func New(b *BulletData, player entity.Entiter, direction int, lifespan time.Time, spawnCount int) (*Bullet, error) {
	name := "base"
	layer, err := library.Layer(b.SpriteName, b.LayerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", name)
	}

	n := &Bullet{
		data:         b,
		player:       player,
		spriteName:   b.SpriteName,
		layerName:    b.LayerName,
		direction:    direction,
		lifespan:     lifespan,
		behaviorType: b.BehaviorType,
		layer:        layer,
		entityID:     entity.NextEntityID(),
		image:        layer.Cells[0].EbitenImage,
		distance:     b.Distance,
		moveSpeed:    b.MoveSpeed,
		isImmortal:   b.IsImmortal,
		spawnCount:   spawnCount,
	}
	library.AudioPlay(b.SpriteName)
	n.x = player.X() + b.OffsetX // + float64(player.SWidth()/4)
	n.y = player.Y() + b.OffsetY // + float64(player.SHeight()/4)

	n.spawnX = n.x
	n.spawnY = n.y

	if b.BehaviorType == BehaviorCircle {
		n.rotation = float64(spawnCount) * 0.2
	}

	err = n.SetAnimation("walk")
	if err != nil {
		return nil, fmt.Errorf("SetAnimation %s: %w", "walk", err)
	}

	bullets = append(bullets, n)
	err = entity.Register(n)
	if err != nil {
		return nil, fmt.Errorf("entity.Register: %w", err)
	}
	return n, nil
}

func (n *Bullet) IsHit(x, y float64) bool {
	if n.IsDead() {
		return false
	}
	return n.image.At(int(x), int(y)).(color.RGBA).A > 0
}

func (n *Bullet) Draw(screen *ebiten.Image) error {
	n.animationStep()
	if len(n.layer.Cells) <= int(n.animation.index) {
		return fmt.Errorf("animationIndex %d is out of bounds for body cells %d", n.animation.index, len(n.layer.Cells))
	}
	c := n.layer.Cells[int(n.animation.index)]

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(n.layer.SpriteWidth/2), -float64(n.layer.SpriteHeight/2))
	op.GeoM.Translate(float64(c.PositionX), float64(c.PositionY))
	op.GeoM.Translate(n.x, n.y)
	op.GeoM.Translate(camera.X, camera.Y)
	op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())

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
	x -= float64(n.SWidth() / 2)
	y += float64(n.SHeight() / 2)

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

func (n *Bullet) SetPosition(x, y float64) {
	n.x, n.y = x, y
}

func (n *Bullet) Position() (float64, float64) {
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

func (n *Bullet) X() float64 {
	return n.x
}

func (n *Bullet) Y() float64 {
	return n.y
}
