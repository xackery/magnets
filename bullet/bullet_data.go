package bullet

import "time"

const (
	BulletNone = iota
	BulletPlain
)

const (
	BehaviorLinear = iota
	BehaviorBoomerang
	BehaviorCircle
)

type BulletData struct {
	Damage       int       // Damage is the amount of damage dealt by the object
	SpriteName   string    // Sprite is the file name
	LayerName    string    // Layer is the layer on aseprite to use
	BehaviorType int       // Behavior handles movement style
	Distance     float64   // Distance is how far a bullet will travel
	IsImmortal   bool      // Lifespam is unlimited if this is true
	Lifespan     time.Time // Lifespan is how long a bullet will be alive for, if IsImmortal is false
	MoveSpeed    float64   // MoveSpeed is how fast a bullet moves at
}
