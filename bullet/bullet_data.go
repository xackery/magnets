package bullet

import "time"

const (
	BulletNone = iota
	BulletPlain
)

const (
	BehaviorNone = iota
	BehaviorLinear
	BehaviorBoomerang
	BehaviorCircle
	BehaviorLasso
	BehaviorWave
	BehaviorUp
	BehaviorKnockback
)

type BulletData struct {
	Damage           int       // Damage is the amount of damage dealt by the object
	SpriteName       string    // Sprite is the file name
	LayerName        string    // Layer is the layer on aseprite to use
	BehaviorType     int       // Behavior handles movement style
	Distance         float64   // Distance is how far a bullet will travel
	IsImmortal       bool      // Lifespam is unlimited if this is true
	Lifespan         time.Time // Lifespan is how long a bullet will be alive for, if IsImmortal is false
	MoveSpeed        float64   // MoveSpeed is how fast a bullet moves at
	SourceWeaponType int       // SourceWeaponType is the type of weapon that is the source of the bullet
	OffsetX          float64   // OffsetX is the bullet's offset on spawn
	OffsetY          float64   // OffsetY is the bullet's offset on spawn
}
