package bullet

const (
	BulletNone = iota
	BulletPlain
)

var (
	bulletTypes = map[int]*bulletData{}
)

type bulletData struct {
	damage     int
	spriteName string
	layerName  string
}

func init() {
	bulletTypes = make(map[int]*bulletData)
	bulletTypes[BulletPlain] = &bulletData{
		damage:     1,
		spriteName: "bullet",
		layerName:  "default",
	}
}
