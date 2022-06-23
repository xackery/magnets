package bullet

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/npc"
)

var (
	bullets []*Bullet
)

func init() {
	global.SubscribeOnResolutionChange(onResolutionChange)
}

func MoveToFront(bullet *Bullet) {
	index := -1
	isFound := false
	for i, it := range bullets {
		if bullet != it {
			continue
		}
		index = i
		isFound = true
		break
	}
	if !isFound {
		return
	}

	bullets = append(bullets[:index], bullets[index+1:]...)
	bullets = append(bullets, bullet)
}

func At(x, y float64) *Bullet {
	for _, p := range bullets {
		if !p.IsHit(x, y) {
			continue
		}
		return p
	}
	return nil
}

func Draw(screen *ebiten.Image) error {
	for _, b := range bullets {
		if b.IsDead() {
			continue
		}
		err := b.Draw(screen)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}
	return nil
}

func Update() {
	for _, b := range bullets {
		if b.IsDead() {
			continue
		}
		b.bulletMove()
	}
}

func HitUpdate() {
	isCleanupNeeded := false

	for _, b := range bullets {
		if b.IsDead() {
			isCleanupNeeded = true
			continue
		}

		for x := 0; x < int(b.layer.SpriteWidth); x++ {
			for y := 0; y < int(b.layer.SpriteHeight); y++ {
				if b.image.At(x, y).(color.RGBA).A == 0 {
					continue
				}
				n := npc.At(b.x+float64(x), b.y+float64(y))
				if n == nil {
					continue
				}
				n.Damage(b.damage)
				if !b.isImmortal {
					b.isDead = true
				}
				return
			}
		}

	}

	if isCleanupNeeded {
		cleanupDead()
	}
}

func Key(index int) *Bullet {
	for i, n := range bullets {
		if i+1 == index {
			return n
		}
	}
	return nil
}

func onResolutionChange() {
	/*	for _, n := range bullets {
		n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	}*/
}

func Bullets() []*Bullet {
	return bullets
}

func HID(hid string) *Bullet {
	for _, n := range bullets {
		if n.HID() == hid {
			return n
		}
	}
	return nil
}

func Clear() {
	bullets = []*Bullet{}
}

func cleanupDead() {
	newBullets := []*Bullet{}

	for _, p := range bullets {
		if p.IsDead() {
			continue
		}
		newBullets = append(newBullets, p)
	}
	bullets = newBullets
}
