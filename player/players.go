package player

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/item"
	"github.com/xackery/magnets/library"
	"github.com/xackery/magnets/npc"
	"github.com/xackery/magnets/ui/life"
	"github.com/xackery/magnets/weapon"
)

var (
	players []*Player
)

func init() {
	global.SubscribeOnResolutionChange(onResolutionChange)
}

func MoveToFront(player *Player) {
	index := -1
	isFound := false
	for i, it := range players {
		if player != it {
			continue
		}
		index = i
		isFound = true
		break
	}
	if !isFound {
		return
	}

	players = append(players[:index], players[index+1:]...)
	players = append(players, player)
}

func At(x, y float64) *Player {
	for _, p := range players {
		if !p.IsHit(x, y) {
			continue
		}
		return p
	}
	return nil
}

func Draw(screen *ebiten.Image) error {
	var err error
	for _, p := range players {
		err = p.Draw(screen)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}

	return nil
}

func Update() {
	for _, p := range players {
		p.Update()
	}
}

func HitUpdate() {
	isCleanupNeeded := false
	isHit := false

	for _, p := range players {
		if p.IsDead() {
			isCleanupNeeded = true
			continue
		}

		for x := 0; x < int(p.layer.SpriteWidth); x++ {
			for y := 0; y < int(p.layer.SpriteHeight); y++ {
				if p.image.At(x, y).(color.RGBA).A == 0 {
					continue
				}
				n := npc.At(p.x+float64(x), p.y+float64(y))
				if n != nil {
					p.Damage(1)
					isHit = true
					if p.IsDead() {
						isCleanupNeeded = true
					}
				}
				i := item.Pickup(p.x+float64(x), p.y+float64(y))
				if i != nil {
					if i.Data.SpriteName == "rupee" {
						p.addExp(i.Data.Value)
					}
					if i.Data.SpriteName == "heart" {

						if p.hp < p.maxHP {
							library.AudioPlay("heal")
							p.hp++
							life.SetHP(p.hp)
						}
					}
				}
			}
			if isHit {
				break
			}
		}
		if isHit {
			break
		}
	}

	if isCleanupNeeded {
		cleanupDead()
	}
}

func Key(index int) *Player {
	for i, n := range players {
		if i+1 == index {
			return n
		}
	}
	return nil
}

func onResolutionChange() {
	/*for _, n := range players {
		n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	}*/

	camera.X = float64(global.ScreenWidth() / 4)
	camera.Y = float64(global.ScreenHeight() / 4)
	if len(players) > 0 {
		camera.X -= players[0].x
		camera.Y -= players[0].y
	}

	for _, n := range players {
		n.expBar.SetWidth(float64(global.ScreenWidth()))
		n.expBar.SetHeight(16)
	}

}

func Players() []*Player {
	return players
}

func HID(hid string) *Player {
	for _, n := range players {
		if n.HID() == hid {
			return n
		}
	}
	return nil
}

func Clear() {
	players = []*Player{}
}

func X() float64 {
	if len(players) == 0 {
		return 0
	}
	return players[0].x
}

func Y() float64 {
	if len(players) == 0 {
		return 0
	}
	return players[0].y
}

func cleanupDead() {
	newPlayers := []*Player{}

	for _, p := range players {
		if p.IsDead() {
			continue
		}
		newPlayers = append(newPlayers, p)
	}
	players = newPlayers
}

func WeaponAdd(weapon *weapon.Weapon) {
	if len(players) == 0 {
		return
	}
	n := players[0]
	n.weapons[weapon.WeaponType] = weapon
}

func IsDead() bool {
	if len(players) == 0 {
		return true
	}
	return players[0].IsDead()
}
