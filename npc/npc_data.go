package npc

const (
	NpcNone = iota
	NpcBat
	NpcCloud
	NpcFlower
	NpcAseprite
	NpcPot
	NpcKnight
	NpcStump
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
		MoveSpeed: 0.2,
		Sprite: &SpriteData{
			spriteName: "flying",
			layerName:  "bat",
		},
	}
	npcTypes[NpcCloud] = &NpcData{
		MaxHP:     10,
		MoveSpeed: 0.2,
		Sprite: &SpriteData{
			spriteName: "flying",
			layerName:  "cloud",
		},
	}
	npcTypes[NpcFlower] = &NpcData{
		MaxHP:     15,
		MoveSpeed: 0.2,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "flower",
		},
	}

	npcTypes[NpcAseprite] = &NpcData{
		MaxHP:     20,
		MoveSpeed: 0.5,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "aseprite",
		},
	}
	npcTypes[NpcPot] = &NpcData{
		MaxHP:     10,
		MoveSpeed: 0.15,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "pot",
		},
	}
	npcTypes[NpcKnight] = &NpcData{
		MaxHP:     150,
		MoveSpeed: 0.3,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "knight",
		},
	}
	npcTypes[NpcStump] = &NpcData{
		MaxHP:     30,
		MoveSpeed: 0.1,
		Sprite: &SpriteData{
			spriteName: "land",
			layerName:  "stump",
		},
	}
}
