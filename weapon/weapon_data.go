package weapon

import (
	"time"

	"github.com/xackery/magnets/bullet"
)

const (
	WeaponNone = iota
	WeaponBoomerang
	WeaponCrystal
	WeaponBoot
	WeaponShovel
	WeaponHammer
	WeaponMagneticGloves
	WeaponSpear
	WeaponShuriken
	// leave this on bottom
	WeaponMax // max value for weapon types
)

var (
	weaponTypes = map[int]*WeaponData{}
)

type WeaponData struct {
	name       string
	Delay      time.Duration
	MaxBullets int
	Bullet     *bullet.BulletData
	Icon       *SpriteData
}

type SpriteData struct {
	SpriteName string
	LayerName  string
}

func init() {
	weaponTypes = make(map[int]*WeaponData)
	weaponTypes[WeaponBoomerang] = &WeaponData{
		name:  "Boomerang",
		Delay: 2500 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorBoomerang,
			SourceWeaponType: WeaponBoomerang,
			Damage:           10,
			SpriteName:       "boomerang",
			LayerName:        "base",
			Distance:         150,
			MoveSpeed:        2,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "boomerang",
		},
	}

	weaponTypes[WeaponCrystal] = &WeaponData{
		name:       "Crystal",
		Delay:      6000 * time.Millisecond,
		MaxBullets: 1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorCircle,
			SourceWeaponType: WeaponCrystal,
			Damage:           10,
			IsImmortal:       true,
			SpriteName:       "crystal",
			LayerName:        "crystal",
			Distance:         50,
			MoveSpeed:        2,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "crystal",
		},
	}

	weaponTypes[WeaponMagneticGloves] = &WeaponData{
		name:       "Magnetic Gloves",
		Delay:      9999 * time.Millisecond,
		MaxBullets: -1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorNone,
			SourceWeaponType: WeaponMagneticGloves,
			IsImmortal:       true,
			SpriteName:       "gloves",
			LayerName:        "gloves",
			Distance:         50,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "gloves",
		},
	}

	weaponTypes[WeaponHammer] = &WeaponData{
		name:       "Hammer",
		Delay:      1000 * time.Millisecond,
		MaxBullets: 1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorLinear,
			SourceWeaponType: WeaponHammer,
			Damage:           20,
			IsImmortal:       true,
			SpriteName:       "hammer",
			LayerName:        "base",
			Distance:         1,
			MoveSpeed:        0,
			OffsetX:          20,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "hammer",
		},
	}

	weaponTypes[WeaponShovel] = &WeaponData{
		name:       "Shovel",
		Delay:      6000 * time.Millisecond,
		MaxBullets: 1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorLasso,
			SourceWeaponType: WeaponShovel,
			Damage:           10,
			IsImmortal:       true,
			SpriteName:       "shovel",
			LayerName:        "base",
			Distance:         50,
			MoveSpeed:        2,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "shovel",
		},
	}

	weaponTypes[WeaponShuriken] = &WeaponData{
		name:  "Shuriken",
		Delay: 900 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorWave,
			SourceWeaponType: WeaponShuriken,
			Damage:           1,
			SpriteName:       "arrow",
			LayerName:        "shuriken",
			Distance:         300,
			MoveSpeed:        4,
		},
		Icon: &SpriteData{
			SpriteName: "arrow",
			LayerName:  "shuriken",
		},
	}

	weaponTypes[WeaponSpear] = &WeaponData{
		name:  "Spear",
		Delay: 1200 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorWave,
			SourceWeaponType: WeaponSpear,
			Damage:           1,
			SpriteName:       "arrow",
			LayerName:        "spear",
			Distance:         300,
			MoveSpeed:        4,
		},
		Icon: &SpriteData{
			SpriteName: "arrow",
			LayerName:  "spear",
		},
	}

	weaponTypes[WeaponBoot] = &WeaponData{
		name:       "Boot",
		Delay:      9999 * time.Millisecond,
		MaxBullets: -1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorNone,
			SourceWeaponType: WeaponBoot,
			IsImmortal:       true,
			SpriteName:       "boot",
			LayerName:        "base",
			Distance:         50,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "boot",
		},
	}
}

func (w *WeaponData) Name() string {
	return w.name
}

func WeaponInfo(weaponType int) *WeaponData {
	return weaponTypes[weaponType]
}
