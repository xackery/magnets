package global

var (
	isPaused bool
)

func IsPaused() bool {
	return isPaused
}

func SetIsPaused(value bool) {
	isPaused = value
}
