package npc

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/magnets/global"
)

var (
	npcs []*Npc
)

func Init() {
	global.SubscribeOnResolutionChange(onResolutionChange)
}

func MoveToFront(npc *Npc) {
	index := -1
	isFound := false
	for i, it := range npcs {
		if npc != it {
			continue
		}
		index = i
		isFound = true
		break
	}
	if !isFound {
		return
	}

	npcs = append(npcs[:index], npcs[index+1:]...)
	npcs = append(npcs, npc)
}

func At(x, y int) *Npc {
	for _, p := range npcs {
		if !p.IsHit(x, y) {
			continue
		}
		return p
	}
	return nil
}

func Draw(screen *ebiten.Image) error {
	var err error
	for _, p := range npcs {
		err = p.Draw(screen)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}

	return nil
}

func Key(index int) *Npc {
	for i, n := range npcs {
		if i+1 == index {
			return n
		}
	}
	return nil
}

func onResolutionChange() {
	for _, n := range npcs {
		if n.dragX != 0 || n.dragY != 0 {
			continue
		}
		n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	}
}

func Npcs() []*Npc {
	return npcs
}

func HID(hid string) *Npc {
	for _, n := range npcs {
		if n.HID() == hid {
			return n
		}
	}
	return nil
}

func Clear() {
	npcs = []*Npc{}
}
