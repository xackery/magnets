package global

const (
	DirectionUp = iota
	DirectionUpRight
	DirectionRight
	DirectionDownRight
	DirectionDown
	DirectionDownLeft
	DirectionLeft
	DirectionUpLeft
)

func IsDirectionLeft(direction int) bool {
	if direction == DirectionDownLeft {
		return true
	}
	if direction == DirectionLeft {
		return true
	}
	if direction == DirectionUpLeft {
		return true
	}
	return false
}

func IsDirectionRight(direction int) bool {
	if direction == DirectionDownRight {
		return true
	}
	if direction == DirectionRight {
		return true
	}
	if direction == DirectionUpRight {
		return true
	}
	return false
}

func IsDirectionDown(direction int) bool {
	if direction == DirectionDownRight {
		return true
	}
	if direction == DirectionDown {
		return true
	}
	if direction == DirectionDownLeft {
		return true
	}
	return false
}

func IsDirectionUp(direction int) bool {
	if direction == DirectionUpRight {
		return true
	}
	if direction == DirectionUp {
		return true
	}
	if direction == DirectionUpLeft {
		return true
	}
	return false
}
