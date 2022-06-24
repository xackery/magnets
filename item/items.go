package item

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

var (
	items []*Item
)

func init() {
	global.SubscribeOnResolutionChange(onResolutionChange)
}

func MoveToFront(item *Item) {
	index := -1
	isFound := false
	for i, it := range items {
		if item != it {
			continue
		}
		index = i
		isFound = true
		break
	}
	if !isFound {
		return
	}

	items = append(items[:index], items[index+1:]...)
	items = append(items, item)
}

func At(x, y float64) *Item {
	for _, p := range items {
		if !p.IsHit(x, y) {
			continue
		}
		return p
	}
	return nil
}

func Draw(screen *ebiten.Image) error {
	for _, b := range items {
		if b.IsDead() {
			continue
		}
		err := b.Draw(screen)
		if err != nil {
			return fmt.Errorf("draw: %w", err)
		}
	}
	return nil
}

func Update() {
	for _, b := range items {
		if b.IsDead() {
			continue
		}
		b.move()
	}
}

func HitUpdate() {

}

func Key(index int) *Item {
	for i, n := range items {
		if i+1 == index {
			return n
		}
	}
	return nil
}

func onResolutionChange() {
	/*	for _, n := range items {
		n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	}*/
}

func Items() []*Item {
	return items
}

func HID(hid string) *Item {
	for _, n := range items {
		if n.HID() == hid {
			return n
		}
	}
	return nil
}

func Clear() {
	items = []*Item{}
}

func Pickup(x, y float64) *Item {
	for _, p := range items {
		if !p.IsHit(x, y) {
			continue
		}
		library.AudioPlay("rupee")
		p.isDead = true
		cleanupDead()
		return p
	}
	return nil
}

func cleanupDead() {
	newItems := []*Item{}

	for _, p := range items {
		if p.IsDead() {
			continue
		}
		newItems = append(newItems, p)
	}
	items = newItems
}
