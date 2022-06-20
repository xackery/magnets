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

func AnchorPosition(anchor int, xOffset, yOffset float64) (float64, float64) {
	switch anchor {
	case AnchorNone:
		return xOffset, yOffset
	case AnchorTopLeft:
		return xOffset, yOffset
	case AnchorTop:
		return xOffset, 0
	case AnchorTopRight:
		return float64(ScreenWidth()) + xOffset, yOffset
	case AnchorRight:
		return float64(ScreenWidth()) + xOffset, yOffset
	case AnchorBottomRight:
		return float64(ScreenWidth()) + xOffset, float64(ScreenHeight()) + yOffset
	case AnchorBottom:
		return xOffset, float64(ScreenHeight()) + yOffset
	case AnchorBottomLeft:
		return xOffset, float64(ScreenHeight()) + yOffset
	}
	return xOffset, yOffset
}
