package weapon

import (
	"time"

	"github.com/xackery/magnets/bullet"
)

const (
	WeaponNone = iota
	WeaponBoomerang
)

var (
	weaponTypes = map[int]*WeaponData{}
)

type WeaponData struct {
	Damage int
	Delay  time.Duration
	Bullet *bullet.BulletData
	Icon   *SpriteData
}

type SpriteData struct {
	spriteName string
	layerName  string
}

func init() {
	weaponTypes = make(map[int]*WeaponData)
	weaponTypes[WeaponBoomerang] = &WeaponData{
		Damage: 1,
		Delay:  500 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType: bullet.BehaviorBoomerang,
			SpriteName:   "bullet",
			LayerName:    "default",
			Distance:     300,
			MoveSpeed:    4,
		},
		Icon: &SpriteData{
			spriteName: "item",
			layerName:  "boomerang",
		},
	}
}
