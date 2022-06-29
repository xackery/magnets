package global

var (
	isLevelUp bool
)

func IsLevelUp() bool {
	return isLevelUp
}

func SetIsLevelUp(value bool) {
	isLevelUp = value
}
