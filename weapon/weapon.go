package weapon

import (
	"fmt"
	"time"

	"github.com/xackery/magnets/bullet"
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

// Name returns a string form of a weapon's name, based on provided weaponType
func Name(weaponType int) string {
	weaponData, ok := weaponTypes[weaponType]
	if !ok {
		return "Unknown"
	}
	return weaponData.name
}
