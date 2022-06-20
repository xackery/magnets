package npc

const (
	NpcNone = iota
	NpcBat
	NpcCloud
	NpcFlower
	NpcAseprite
	NpcPot
	NpcKnight
)

var (
	npcTypes = map[int]*NpcData{}
)

type NpcData struct {
	MaxHP     int
	Sprite    *SpriteData
	MoveSpeed float64
}

type SpriteData struct {
	spriteName string
	layerName  string
}

func init() {
	npcTypes = make(map[int]*NpcData)
	npcTypes[NpcBat] = &NpcData{
		MaxHP:     5,
		MoveSpeed: 0.1,
		Sprite: &SpriteData{
			spriteName: "flying",
			layerName:  "bat",
		},
	}
	npcTypes[NpcCloud] = &NpcData{
		MaxHP:     5,
		MoveSpeed: 0.1,
		Sprite: &SpriteData{
			spriteName: "flying",
			layerName:  "cloud",
		},
	}
	npcTypes[NpcFlower] = &NpcData{
		MaxHP:     5,
		MoveSpeed: 0.1,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "flower",
		},
	}

	npcTypes[NpcAseprite] = &NpcData{
		MaxHP:     5,
		MoveSpeed: 0.1,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "aseprite",
		},
	}
	npcTypes[NpcPot] = &NpcData{
		MaxHP:     5,
		MoveSpeed: 0.1,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "pot",
		},
	}
	npcTypes[NpcKnight] = &NpcData{
		MaxHP:     5,
		MoveSpeed: 0.1,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "knight",
		},
	}
}
