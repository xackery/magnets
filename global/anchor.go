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

func AnchorPosition(anchor, xOffset, yOffset int) (int, int) {
	switch anchor {
	case AnchorNone:
		return xOffset, yOffset
	case AnchorTopLeft:
		return xOffset, yOffset
	case AnchorTop:
		return xOffset, 0
	case AnchorTopRight:
		return ScreenWidth() + xOffset, yOffset
	case AnchorRight:
		return ScreenWidth() + xOffset, yOffset
	case AnchorBottomRight:
		return ScreenWidth() + xOffset, ScreenHeight() + yOffset
	case AnchorBottom:
		return xOffset, ScreenHeight() + yOffset
	case AnchorBottomLeft:
		return xOffset, ScreenHeight() + yOffset
	}
	return xOffset, yOffset
}
