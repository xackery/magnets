package bullet

import "time"

const (
	BulletNone = iota
	BulletPlain
)

const (
	BehaviorLinear = iota
	BehaviorBoomerang
)

type BulletData struct {
	Damage       int       // Damage is the amount of damage dealt by the object
	SpriteName   string    // Sprite is the file name
	LayerName    string    // Layer is the layer on aseprite to use
	BehaviorType int       // Behavior handles movement style
	Distance     float32   // Distance is how far a bullet will travel
	Lifespan     time.Time // Lifespan is how long a bullet will be alive for
	MoveSpeed    float32   // MoveSpeed is how fast a bullet moves at
}
