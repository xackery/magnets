package player

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xackery/magnets/bullet"
	"github.com/xackery/magnets/weapon"
)

func (n *Player) WeaponAdd(weapon *weapon.Weapon) {
	_, ok := n.weapons[weapon.WeaponType]
	if ok {
		return
	}
	n.weapons[weapon.WeaponType] = weapon
	log.Debug().Msgf("player equipped weapon type %d", weapon.WeaponType)
}

func (n *Player) weaponUpdate() {
	for _, weapon := range n.weapons {
		if weapon.Data.Bullet.IsImmortal && len(weapon.Bullets) > 0 {
			continue
		} else if time.Now().Before(weapon.Cooldown) {
			continue
		}

		weapon.Shoot()
		b, err := bullet.New(weapon.Data.Bullet, n, n.direction)
		if err != nil {
			fmt.Println("bullet.New:", err)
		}
		weapon.Bullets = append(weapon.Bullets, b)
	}
}
