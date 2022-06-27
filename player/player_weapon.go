package player

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/xackery/magnets/bullet"
	"github.com/xackery/magnets/ui/equipment"
	"github.com/xackery/magnets/weapon"
)

func (n *Player) weaponAdd(w *weapon.Weapon) error {
	_, ok := n.weapons[w.WeaponType]
	if ok {
		return fmt.Errorf("weapon %d not found", w.WeaponType)
	}
	n.weapons[w.WeaponType] = w
	_, err := equipment.New(w.Data)
	if err != nil {
		return fmt.Errorf("equipment.New %s: %w", weapon.Name(w.WeaponType), err)
	}
	log.Debug().Msgf("player equipped %s (%d)", weapon.Name(w.WeaponType), w.WeaponType)
	return nil
}

func (n *Player) weaponUpdate() {
	for _, w := range n.weapons {
		if w.Data.MaxBullets == -1 {
			continue
		}
		i := 0
		for _, b := range w.Bullets {
			if b == nil {
				continue
			}
			if b.IsDead() {
				continue
			}
			w.Bullets[i] = b
			i++
		}
		for j := i; j < len(w.Bullets); j++ {
			w.Bullets[j] = nil
		}
		w.Bullets = w.Bullets[:i]

		if len(w.Bullets) >= n.weaponMaxBullets(w.WeaponType) {
			continue
		}

		if time.Now().Before(w.Cooldown) {
			continue
		}

		w.Cooldown = time.Now().Add(n.weaponDelay(w.WeaponType))

		n.weaponBulletSpawn(w)
	}
}

func (n *Player) hasWeapon(weaponType int) bool {
	_, ok := n.weapons[weaponType]
	return ok
}

func (n *Player) weaponUpgrade(weaponType int) {
	n.weapons[weaponType].Level++
	if n.weapons[weaponType].Level > 9 {
		n.weapons[weaponType].Level = 9
	}
	log.Debug().Msgf("weapon %s (%d) is now level %d", weapon.Name(weaponType), weaponType, n.weapons[weaponType].Level)
}

func (n *Player) weaponDelay(weaponType int) time.Duration {
	w, ok := n.weapons[weaponType]
	if !ok {
		return 1 * time.Second
	}

	delay := w.Data.Delay
	if w.Level == 1 {
		return delay
	}

	if weaponType == weapon.WeaponBoomerang && w.Level >= 2 {
		delay /= time.Duration(w.Level)
	}

	return delay

}

func (n *Player) weaponMaxBullets(weaponType int) int {
	w, ok := n.weapons[weaponType]
	if !ok {
		return 1
	}
	if w.Data.MaxBullets == 0 {
		return 999
	}

	if weaponType == weapon.WeaponCrystal && w.Level > 1 {
		return w.Data.MaxBullets + w.Level
	}
	return w.Data.MaxBullets
}

func (n *Player) weaponBulletLifespan(weaponType int) time.Duration {
	lifespan := 10 * time.Second
	w, ok := n.weapons[weaponType]
	if !ok {
		return lifespan
	}
	if w.WeaponType == weapon.WeaponCrystal {
		lifespan = 1700 * time.Millisecond
	}
	if w.WeaponType == weapon.WeaponHammer {
		lifespan = 300 * time.Millisecond
	}
	return lifespan
}

func (n *Player) weaponBulletSpawn(w *weapon.Weapon) {
	count := 1

	if w.WeaponType == weapon.WeaponCrystal && w.Level > 1 {
		count = w.Level
	}

	for i := 0; i < count; i++ {
		b, err := bullet.New(w.Data.Bullet, n, n.direction, time.Now().Add(n.weaponBulletLifespan(w.WeaponType)), i)
		if err != nil {
			fmt.Println("bullet.New:", err)
		}
		w.Bullets = append(w.Bullets, b)
	}
}

func (n *Player) weaponLevel(weaponType int) int {
	w, ok := n.weapons[weaponType]
	if !ok {
		return 0
	}
	return w.Level
}
