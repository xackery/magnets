package global

const (
	AnchorNone = iota
	AnchorTopLeft
	AnchorTop
	AnchorTopRight
	AnchorRight
	AnchorBottomRight
	AnchorBottom
	AnchorBottomLeft
	AnchorLeft
)

func AnchorPosition(anchor int, xOffset, yOffset float32) (float32, float32) {
	switch anchor {
	case AnchorNone:
		return xOffset, yOffset
	case AnchorTopLeft:
		return xOffset, yOffset
	case AnchorTop:
		return xOffset, 0
	case AnchorTopRight:
		return float32(ScreenWidth()) + xOffset, yOffset
	case AnchorRight:
		return float32(ScreenWidth()) + xOffset, yOffset
	case AnchorBottomRight:
		return float32(ScreenWidth()) + xOffset, float32(ScreenHeight()) + yOffset
	case AnchorBottom:
		return xOffset, float32(ScreenHeight()) + yOffset
	case AnchorBottomLeft:
		return xOffset, float32(ScreenHeight()) + yOffset
	}
	return xOffset, yOffset
}
