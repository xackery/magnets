package global

var (
	isGameOver bool
)

func IsGameOver() bool {
	return isGameOver
}

func SetIsGameOver(value bool) {
	isGameOver = value
}
