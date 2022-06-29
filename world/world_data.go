package world

const (
	WorldNone = iota
	WorldGrass
)

var (
	worldTypes = map[int]*WorldData{}
)

type WorldData struct {
	Tilemap *SpriteData
	Icon    *SpriteData
}

type SpriteData struct {
	spriteName string
	layerName  string
}

func init() {
	/*worldTypes = make(map[int]*WorldData)
	worldTypes[WorldGrass] = &WorldData{
		Tilemap: &SpriteData{
			spriteName: "world-1",
			layerName:  "base",
		},
		Icon: &SpriteData{
			spriteName: "world-1",
			layerName:  "base",
		},
	}
	*/
}
