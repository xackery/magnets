package npc

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog/log"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/input"
)

var (
	npcs                []*Npc
	isAIEnabled         bool
	isAIEnabledCooldown time.Time
	spawnerCooldown     time.Time
	spawnMax            int = 50
)

func init() {
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

func At(x, y float64) *Npc {
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

func Update(playerX, playerY float64) {
	if input.IsPressed(ebiten.KeyGraveAccent) && time.Now().After(isAIEnabledCooldown) {
		isAIEnabled = !isAIEnabled
		log.Debug().Msgf("AI is now %t", isAIEnabled)
		isAIEnabledCooldown = time.Now().Add(500 * time.Millisecond)
	}

	for _, p := range npcs {
		p.Update(playerX, playerY)
	}
	if isAIEnabled {
		if time.Now().After(spawnerCooldown) {
			spawnerCooldown = time.Now().Add(3 * time.Second)
			if len(npcs) > spawnMax {
				return
			}
			spawnCount := spawnMax + 1 - len(npcs)
			maxDistance := global.ScreenSmallestDimension() / 2
			minDistance := global.ScreenSmallestDimension() / 2.5
			for i := 0; i < spawnCount; i++ {

				theta := rand.Float64() * (math.Pi * 2) * math.Pi

				distance := (rand.Float64() * (maxDistance - minDistance)) + minDistance
				//fmt.Println("distance", distance, global.ScreenSmallestDimension())
				/*
					rX := -100 + rand.Float64()*(100-200)
					rY := -100 + rand.Float64()*(100-200)*/

				New(NpcBat, playerY+math.Sin(theta)*distance, playerX+math.Cos(theta)*distance)
			}
			log.Debug().Msgf("spawned %d", spawnCount)
		}
	}
}

func HitUpdate() {
	isCleanupNeeded := false

	for _, p := range npcs {
		if p.IsDead() {
			isCleanupNeeded = true
			continue
		}
	}

	if isCleanupNeeded {
		cleanupDead()
	}
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
	/*for _, n := range npcs {
		n.SetPosition(global.AnchorPosition(n.anchor, n.xOffset, n.yOffset))
	}*/
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

func cleanupDead() {
	newNpcs := []*Npc{}

	for _, p := range npcs {
		if p.IsDead() {
			continue
		}
		newNpcs = append(newNpcs, p)
	}
	npcs = newNpcs
}
