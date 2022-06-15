package global

import "github.com/rs/zerolog/log"

var (
	screenWidth                   int
	screenHeight                  int
	screenScaleX                  float64 = 1
	screenScaleY                  float64 = 1
	onResolutionChangeSubscribers []func()
)

func ScreenWidth() int {
	return screenWidth
}

func ScreenHeight() int {
	return screenHeight
}

func ScreenIsLandscape() bool {
	return screenWidth > screenHeight
}

func ScreenOnLayoutChange(x, y int) {
	if (screenWidth == x && screenHeight == y) || (x == 0 && y == 0) {
		return
	}
	log.Debug().Msgf("layout changed to %dx%d", x, y)
	screenWidth = x
	screenHeight = y
	for _, event := range onResolutionChangeSubscribers {
		event()
	}
}

func ScreenScaleChange(x float64, y float64) {
	screenScaleX = x
	screenScaleY = y
}

func ScreenScaleX() float64 {
	return screenScaleX
}

func ScreenScaleY() float64 {
	return screenScaleY
}

func SubscribeOnResolutionChange(event func()) {
	onResolutionChangeSubscribers = append(onResolutionChangeSubscribers, event)
}
