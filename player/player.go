package player

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/entity"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/input"
	"github.com/xackery/magnets/library"
	"github.com/xackery/magnets/weapon"
)

type Player struct {
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
	isLastLeft bool
	animation  animation
	direction  int
	weapons    map[int]*weapon.Weapon
}

type animation struct {
	tag              *aseprite.Tag
	delay            time.Time
	index            int16
	isPingPongToggle bool
}

func New(spriteName string, layerName string) (*Player, error) {
	name := "base"
	layer, err := library.Layer(spriteName, layerName)
	if err != nil {
		return nil, fmt.Errorf("library.Layer: %w", err)
	}
	if len(layer.Cells) < 1 {
		return nil, fmt.Errorf("no cells found on layer %s", name)
	}

	n := &Player{
		spriteName: spriteName,
		layerName:  layerName,
		hp:         1,
		maxHP:      1,
		direction:  global.DirectionRight,
		layer:      layer,
		entityID:   entity.NextEntityID(),
		image:      layer.Cells[0].EbitenImage,
		weapons:    make(map[int]*weapon.Weapon),
	}
	global.Player = n

	err = n.SetAnimation("left")
	if err != nil {
		return nil, fmt.Errorf("SetAnimation %s: %w", "left", err)
	}

	players = append(players, n)
	err = entity.Register(n)
	if err != nil {
		return nil, fmt.Errorf("entity.Register: %w", err)
	}
	return n, nil
}

func (n *Player) IsHit(x, y float64) bool {
	if n.IsDead() {
		return false
	}
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

func (n *Player) Draw(screen *ebiten.Image) error {

	if len(n.layer.Cells) <= int(n.animation.index) {
		return fmt.Errorf("animationIndex %d is out of bounds for body cells %d", n.animation.index, len(n.layer.Cells))
	}
	c := n.layer.Cells[int(n.animation.index)]

	op := &ebiten.DrawImageOptions{}
	if n.isLastLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(n.layer.SpriteWidth/2), 0)
	}
	op.GeoM.Translate(-float64(n.layer.SpriteWidth/2), -float64(n.layer.SpriteHeight/2))
	op.GeoM.Translate(float64(c.PositionX), float64(c.PositionY))
	op.GeoM.Translate(float64(camera.X), float64(camera.Y))
	op.GeoM.Translate(float64(n.x), float64(n.y))
	op.GeoM.Scale(global.ScreenScaleX(), global.ScreenScaleY())

	//op.GeoM.Translate(float64(global.ScreenWidth())/2, float64(global.ScreenHeight())/2)
	//op.GeoM.Translate(float64(n.x), float64(n.y))

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
	//x := n.x
	//y := n.y + 30
	//x -= float64(n.SWidth() / 2)
	//y += float64(n.SHeight() / 2)

	return nil
}

// SetAnimation sets the animation of the player
func (n *Player) SetAnimation(name string) error {
	name = strings.ToLower(name)
	tag, err := library.Tag(n.spriteName, name)
	if err != nil {
		return fmt.Errorf("tag: %w", err)
	}
	n.animation.tag = tag
	n.animation.index = tag.From
	return nil
}

func (n *Player) animationStep() {
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

func (n *Player) Update() {

	if input.IsPressed(ebiten.Key1) {
		w, err := weapon.New(weapon.WeaponArrow)
		if err != nil {
			fmt.Println("weapon new:", err)
			os.Exit(1)
		}
		n.WeaponAdd(w)
	}

	isMoving := input.IsPlayerMoving()
	moveSpeed := float64(0.5)
	if isMoving {
		n.direction = input.PlayerDirection
		n.animationStep()
		delta := float64(moveSpeed)
		if global.IsDirectionLeft(n.direction) {
			/*delta = 0
			for ; delta < moveSpeed; delta += 0.1 {
				if world.IsCollision(n.x-delta, n.y) {
					delta -= 0.1
					break
				}
			}*/

			camera.X += delta
			n.x -= delta
			n.isLastLeft = true
		}
		if global.IsDirectionRight(n.direction) {
			/*delta = 0
			for ; delta < moveSpeed; delta += 0.1 {
				if world.IsCollision(n.x+delta, n.y) {
					delta -= 0.1
					break
				}
			}*/
			camera.X -= delta
			n.x += delta
			n.isLastLeft = false
		}

		if global.IsDirectionDown(n.direction) {
			/*delta = 0
			for ; delta < moveSpeed; delta += 0.1 {
				if world.IsCollision(n.x, n.y+delta) {
					delta -= 0.1
					break
				}
			}*/
			camera.Y -= delta
			n.y += delta
		}

		if global.IsDirectionUp(n.direction) {
			/*delta = 0
			for ; delta < moveSpeed; delta += 0.1 {
				if world.IsCollision(n.x, n.y-delta) {
					delta -= 0.1
					break
				}
			}*/
			camera.Y += delta
			n.y -= delta
		}
	}

	n.weaponUpdate()
}

func (n *Player) SetPosition(x, y float64) {
	n.x, n.y = x, y
}

func (n *Player) Position() (float64, float64) {
	return n.x, n.y
}

func (n *Player) HID() string {
	return n.hid
}

func (n *Player) Damage(damage int) bool {
	n.hp -= damage
	if n.hp < 1 {
		log.Debug().Msgf("player got killed")
		n.hp = 0
		return true
	}
	return false
}

func (n *Player) IsDead() bool {
	return n.hp < 1
}

func (n *Player) EntityID() uint {
	return n.entityID
}

func (n *Player) SWidth() int {
	return int(n.layer.SpriteWidth * uint16(global.ScreenScaleX()))
}

func (n *Player) SHeight() int {
	return int(n.layer.SpriteHeight * uint16(global.ScreenScaleY()))
}

func (n *Player) AnimationAttack() error {
	return nil
}

func (n *Player) AnimationGotHit() error {
	return nil
}

func (n *Player) X() float64 {
	return n.x
}

func (n *Player) Y() float64 {
	return n.y
}
