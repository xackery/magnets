package player

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/magnets/global"
)

var (
	players []*Player
)

func Init() {
	global.SubscribeOnResolutionChange(onResolutionChange)
}

func MoveToFront(player *Player) {
	index := -1
	isFound := false
	for i, it := range players {
		if player != it {
			continue
		}
		index = i
		isFound = true
		break
	}
	if !isFound {
		return
	}

	players = append(players[:index], players[index+1:]...)
	players = append(players, player)
}

func At(x, y float32) *Player {
	for _, p := range players {
		if !p.IsHit(x, y) {
			continue
		}
		return p
	}
	return nil
}

func Draw(screen *ebiten.Image) error {
	var err error
	for _, p := range players {
		err = p.Draw(screen)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}

	return nil
}

func Key(index int) *Player {
	for i, n := range players {
		if i+1 == index {
			return n
		}
	}
	return nil
}

func onResolutionChange() {
	for _, n := range players {
		n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	}
}

func Players() []*Player {
	return players
}

func HID(hid string) *Player {
	for _, n := range players {
		if n.HID() == hid {
			return n
		}
	}
	return nil
}

func Clear() {
	players = []*Player{}
}

func X() float32 {
	if len(players) == 0 {
		return 0
	}
	return players[0].x
}

func Y() float32 {
	if len(players) == 0 {
		return 0
	}
	return players[0].y
}
