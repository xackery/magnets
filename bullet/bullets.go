package bullet

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/npc"
)

var (
	bullets []*Bullet
)

func Init() {
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

func At(x, y float32) *Bullet {
	for _, p := range bullets {
		if !p.IsHit(x, y) {
			continue
		}
		return p
	}
	return nil
}

func Draw(screen *ebiten.Image) error {
	var err error
	isCleanupNeeded := false
	for _, p := range bullets {
		if p.IsDead() {
			isCleanupNeeded = true
		}
		err = p.Draw(screen)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}
	if isCleanupNeeded {
		cleanupDead()
	}
	return nil
}

func Update() {
	for _, b := range bullets {
		if b.IsDead() {
			continue
		}
		n := npc.At(b.x-float32(b.layer.SpriteWidth/2), b.y)
		if n != nil {
			n.Damage(b.damage)
			b.isDead = true
			continue
		}
		n = npc.At(b.x+float32(b.layer.SpriteWidth/2), b.y)
		if n != nil {
			n.Damage(b.damage)
			b.isDead = true
			continue
		}
		n = npc.At(b.x, b.y-float32(b.layer.SpriteHeight/2))
		if n != nil {
			n.Damage(b.damage)
			b.isDead = true
			continue
		}
		n = npc.At(b.x, b.y+float32(b.layer.SpriteHeight/2))
		if n != nil {
			n.Damage(b.damage)
			b.isDead = true
			continue
		}

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
	for _, n := range bullets {
		n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	}
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
