package weapon

import (
	"fmt"
	"time"

	"github.com/xackery/magnets/bullet"
)

const (
	None = iota
	Boomerang
)

type Weapon struct {
	Data       *WeaponData
	Cooldown   time.Time
	WeaponType int
	Level      int
	Bullets    []*bullet.Bullet
}

func New(weaponType int) (*Weapon, error) {
	weaponData, ok := weaponTypes[weaponType]
	if !ok {
		return nil, fmt.Errorf("weapon type %d not found", weaponType)
	}
	weapon := &Weapon{
		WeaponType: weaponType,
		Data:       weaponData,
		Level:      1,
	}
	return weapon, nil
}

func (n *Weapon) Delay() time.Duration {
	return n.Data.Delay
}

// Shoot fires a bullet, this sets cooldowns
func (n *Weapon) Shoot() {
	n.Cooldown = time.Now().Add(n.Data.Delay)
}
