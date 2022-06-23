package weapon

import (
	"time"

	"github.com/xackery/magnets/bullet"
)

const (
	WeaponNone = iota
	WeaponBoomerang
	WeaponArrow
	WeaponSword
	WeaponSpear
	WeaponShuriken
)

var (
	weaponTypes = map[int]*WeaponData{}
)

type WeaponData struct {
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
		Delay: 250 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType: bullet.BehaviorBoomerang,
			Damage:       1,
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

	weaponTypes[WeaponArrow] = &WeaponData{
		Delay: 1500 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType: bullet.BehaviorCircle,
			Damage:       10,
			IsImmortal:   true,
			SpriteName:   "crystal",
			LayerName:    "crystal",
			Distance:     20,
			MoveSpeed:    4,
		},
		Icon: &SpriteData{
			spriteName: "crystal",
			layerName:  "crystal",
		},
	}

	weaponTypes[WeaponSword] = &WeaponData{
		Delay: 600 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType: bullet.BehaviorLinear,
			Damage:       1,
			SpriteName:   "arrow",
			LayerName:    "sword",
			Distance:     300,
			MoveSpeed:    4,
		},
		Icon: &SpriteData{
			spriteName: "arrow",
			layerName:  "sword",
		},
	}

	weaponTypes[WeaponShuriken] = &WeaponData{
		Delay: 900 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType: bullet.BehaviorLinear,
			Damage:       1,
			SpriteName:   "arrow",
			LayerName:    "shuriken",
			Distance:     300,
			MoveSpeed:    4,
		},
		Icon: &SpriteData{
			spriteName: "arrow",
			layerName:  "shuriken",
		},
	}

	weaponTypes[WeaponSpear] = &WeaponData{
		Delay: 1200 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType: bullet.BehaviorLinear,
			Damage:       1,
			SpriteName:   "arrow",
			LayerName:    "spear",
			Distance:     300,
			MoveSpeed:    4,
		},
		Icon: &SpriteData{
			spriteName: "arrow",
			layerName:  "spear",
		},
	}
}
