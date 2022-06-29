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
	WeaponMagnet
	WeaponShuriken
	WeaponHeart
	// leave this on bottom
	WeaponMax // max value for weapon types
)

var (
	weaponTypes = map[int]*WeaponData{}
)

type WeaponData struct {
	name        string
	Delay       time.Duration
	MaxBullets  int
	Bullet      *bullet.BulletData
	Description string
	Icon        *SpriteData
}

type SpriteData struct {
	SpriteName string
	LayerName  string
}

func init() {
	weaponTypes = make(map[int]*WeaponData)
	weaponTypes[WeaponBoomerang] = &WeaponData{
		name:        "Boomerang",
		Description: "Fires at an enemy and returns back",
		Delay:       2500 * time.Millisecond,
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
		name:        "Crystal",
		Delay:       6000 * time.Millisecond,
		Description: "Circles around the caster",
		MaxBullets:  1,
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
		name:        "Magnetic Gloves",
		Description: "Knocks nearby enemies back",
		Delay:       10000 * time.Millisecond,
		MaxBullets:  1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorKnockback,
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

	weaponTypes[WeaponMagnet] = &WeaponData{
		name:        "Magnet",
		Description: "Attracts coins to player",
		Delay:       9999 * time.Millisecond,
		MaxBullets:  -1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorNone,
			SourceWeaponType: WeaponMagnet,
			IsImmortal:       true,
			SpriteName:       "magnet",
			LayerName:        "base",
			Distance:         50,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "magnet",
		},
	}

	weaponTypes[WeaponHammer] = &WeaponData{
		name:        "Hammer",
		Description: "Smashes enemies to the east of player",
		Delay:       1000 * time.Millisecond,
		MaxBullets:  1,
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
		name:        "Shovel",
		Description: "Fires at enemies to the north",
		Delay:       4000 * time.Millisecond,
		MaxBullets:  1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorUp,
			SourceWeaponType: WeaponShovel,
			Damage:           10,
			IsImmortal:       false,
			SpriteName:       "shovel",
			LayerName:        "base",
			Distance:         25,
			MoveSpeed:        2,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "shovel",
		},
	}

	weaponTypes[WeaponShuriken] = &WeaponData{
		name:        "Shuriken",
		Description: "Fires a wavy shuriken at enemies",
		Delay:       900 * time.Millisecond,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorWave,
			SourceWeaponType: WeaponShuriken,
			Damage:           4,
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

	weaponTypes[WeaponBoot] = &WeaponData{
		name:        "Boot",
		Description: "Increases movement speed",
		Delay:       9999 * time.Millisecond,
		MaxBullets:  -1,
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

	weaponTypes[WeaponHeart] = &WeaponData{
		name:        "Heart",
		Description: "Increases maximum health",
		Delay:       9999 * time.Millisecond,
		MaxBullets:  -1,
		Bullet: &bullet.BulletData{
			BehaviorType:     bullet.BehaviorNone,
			SourceWeaponType: WeaponHeart,
			IsImmortal:       true,
			SpriteName:       "heart-weapon",
			LayerName:        "base",
			Distance:         50,
		},
		Icon: &SpriteData{
			SpriteName: "icon",
			LayerName:  "heart",
		},
	}
}

func (w *WeaponData) Name() string {
	return w.name
}

func WeaponInfo(weaponType int) *WeaponData {
	return weaponTypes[weaponType]
}
