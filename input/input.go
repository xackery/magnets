package input

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/xackery/magnets/global"
)

var (
	keys            = map[ebiten.Key]*keyEntry{}
	playerIsMoving  bool
	PlayerDirection int
)

type keyEntry struct {
	key           ebiten.Key
	event         func()
	isPressed     bool
	pressStart    time.Time
	lastPressTime time.Duration
}

func init() {
	keys = make(map[ebiten.Key]*keyEntry)
	keys[ebiten.KeyArrowLeft] = &keyEntry{key: ebiten.KeyArrowLeft}
	keys[ebiten.KeyArrowRight] = &keyEntry{key: ebiten.KeyArrowRight}
	keys[ebiten.KeyArrowUp] = &keyEntry{key: ebiten.KeyArrowUp}
	keys[ebiten.KeyArrowDown] = &keyEntry{key: ebiten.KeyArrowDown}
	keys[ebiten.KeyQ] = &keyEntry{key: ebiten.KeyQ}
	keys[ebiten.Key1] = &keyEntry{key: ebiten.Key1}
	keys[ebiten.Key2] = &keyEntry{key: ebiten.Key2}
	keys[ebiten.Key3] = &keyEntry{key: ebiten.Key3}
	keys[ebiten.Key4] = &keyEntry{key: ebiten.Key4}
	keys[ebiten.Key5] = &keyEntry{key: ebiten.Key5}
	keys[ebiten.Key6] = &keyEntry{key: ebiten.Key6}
	keys[ebiten.Key7] = &keyEntry{key: ebiten.Key7}
	keys[ebiten.Key8] = &keyEntry{key: ebiten.Key8}
	keys[ebiten.Key9] = &keyEntry{key: ebiten.Key9}
	keys[ebiten.KeyP] = &keyEntry{key: ebiten.KeyP}
	keys[ebiten.KeyGraveAccent] = &keyEntry{key: ebiten.KeyGraveAccent}
}

func Update() {

	for _, k := range keys {
		if k.isPressed && inpututil.IsKeyJustReleased(k.key) {
			k.isPressed = false
			k.lastPressTime = time.Since(k.pressStart)
			continue
		}
		if !k.isPressed && inpututil.IsKeyJustPressed(k.key) {
			k.isPressed = true
			k.pressStart = time.Now()
			if k.event != nil {
				k.event()
			}
		}
	}
	updatePlayerDirection()
	/*if IsPressed(ebiten.KeyQ) {
		log.Info().Msgf("q pressed, quitting!")
		os.Exit(0)
	}*/
}

func IsPressed(key ebiten.Key) bool {
	k, ok := keys[key]
	if !ok {
		return false
	}
	return k.isPressed
}

func IsPlayerMoving() bool {
	return playerIsMoving
}

func updatePlayerDirection() {
	playerIsMoving = false
	if IsPressed(ebiten.KeyLeft) {
		if IsPressed(ebiten.KeyUp) {
			PlayerDirection = global.DirectionUpLeft
		} else if IsPressed(ebiten.KeyDown) {
			PlayerDirection = global.DirectionDownLeft
		} else {
			PlayerDirection = global.DirectionLeft
		}
		playerIsMoving = true
	}
	if IsPressed(ebiten.KeyRight) {
		if IsPressed(ebiten.KeyUp) {
			PlayerDirection = global.DirectionUpRight
		} else if IsPressed(ebiten.KeyDown) {
			PlayerDirection = global.DirectionDownRight
		} else {
			PlayerDirection = global.DirectionRight
		}
		playerIsMoving = true
	}
	if IsPressed(ebiten.KeyUp) {
		if !IsPressed(ebiten.KeyLeft) && !IsPressed(ebiten.KeyRight) {
			PlayerDirection = global.DirectionUp
		}
		playerIsMoving = true
	}
	if IsPressed(ebiten.KeyDown) {
		if !IsPressed(ebiten.KeyLeft) && !IsPressed(ebiten.KeyRight) {
			PlayerDirection = global.DirectionDown
		}
		playerIsMoving = true
	}

}
