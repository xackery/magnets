package player

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/magnets/bullet"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/weapon"
)

func (n *Player) WeaponAdd(weapon *weapon.Weapon) {
	n.weapons[weapon.WeaponType] = weapon
}

func (n *Player) weaponDraw(screen *ebiten.Image) {
	for _, weapon := range n.weapons {
		if time.Now().Before(weapon.Cooldown) {
			continue
		}
		weapon.Shoot()
		_, err := bullet.New(weapon.Data.Bullet, n, n.direction, global.AnchorTopLeft, n.x, n.y)
		if err != nil {
			fmt.Println("bullet.New:", err)
		}
	}
}
