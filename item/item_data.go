package item

const (
	ItemNone = iota
	ItemRupee
	ItemGreenRupee
	ItemRedRupee
	ItemHeart
)

var (
	itemTypes = make(map[int]*ItemData)
)

type ItemData struct {
	SpriteName string // Sprite is the file name
	LayerName  string // Layer is the layer on aseprite to use
	Value      int    // value, e.g. green rupee is 1
}

func init() {
	itemTypes[ItemRupee] = &ItemData{
		SpriteName: "rupee",
		LayerName:  "base",
		Value:      1,
	}

	itemTypes[ItemGreenRupee] = &ItemData{
		SpriteName: "rupee",
		LayerName:  "green",
		Value:      5,
	}

	itemTypes[ItemRedRupee] = &ItemData{
		SpriteName: "rupee",
		LayerName:  "red",
		Value:      20,
	}

	itemTypes[ItemHeart] = &ItemData{
		SpriteName: "heart",
		LayerName:  "full",
		Value:      1,
	}
}
