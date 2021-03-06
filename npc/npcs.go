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

const (
	maxDistance = float64(600)
	minDistance = float64(200)
)

var (
	npcs                []*Npc
	isAIEnabled         bool = true
	isAIEnabledCooldown time.Time
	spawnerCooldown     time.Time
	step                int
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

func Update() {
	if input.IsPressed(ebiten.KeyGraveAccent) && time.Now().After(isAIEnabledCooldown) {
		isAIEnabled = !isAIEnabled
		log.Debug().Msgf("AI is now %t", isAIEnabled)
		isAIEnabledCooldown = time.Now().Add(500 * time.Millisecond)
	}

	for _, p := range npcs {
		p.Update()
	}
	spawner()
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
	step = 0
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

func spawner() {
	if !isAIEnabled {
		return
	}

	if global.Countdown < 600 && step == 0 {
		spawn(NpcBat, 30)
		step = 1
	}

	if global.Countdown < 550 && step == 1 {
		spawn(NpcBat, 20)
		spawn(NpcCloud, 20)
		step = 2
	}

	if global.Countdown < 500 && step == 2 {
		spawn(NpcKnight, 1)
		spawn(NpcBat, 20)
		spawn(NpcCloud, 20)
		spawn(NpcFlower, 20)
		step = 3
	}

	if global.Countdown < 450 && step == 3 {
		spawn(NpcCloud, 10)
		spawn(NpcFlower, 10)
		step = 4
	}

	if global.Countdown < 400 && step == 4 {
		spawn(NpcPot, 50)
		step = 5
	}

	if global.Countdown < 350 && step == 5 {
		spawn(NpcKnight, 2)
		spawn(NpcCloud, 10)
		spawn(NpcFlower, 10)
		spawn(NpcPot, 10)
		step = 6
	}

	if global.Countdown < 250 && step == 6 {
		spawn(NpcAseprite, 3)
		spawn(NpcPot, 20)
		spawn(NpcFlower, 20)
		spawn(NpcCloud, 10)
		step = 7
	}

	if global.Countdown < 150 && step == 7 {
		spawn(NpcAseprite, 6)
		spawn(NpcPot, 20)
		spawn(NpcFlower, 20)
		spawn(NpcCloud, 10)
		step = 8
	}

	if global.Countdown < 100 && step == 8 {
		spawn(NpcKnight, 20)
		step = 9
	}

	if global.Countdown < 50 && step == 9 {
		spawn(NpcKnight, 20)
		spawn(NpcCloud, 20)
		spawn(NpcFlower, 20)
		spawn(NpcAseprite, 3)
		spawn(NpcPot, 20)
		spawn(NpcFlower, 20)
		spawn(NpcCloud, 10)
		step = 10
	}

	if global.Countdown < 20 && step == 10 {
		spawn(NpcKnight, 20)
		spawn(NpcCloud, 20)
		spawn(NpcFlower, 20)
		spawn(NpcAseprite, 3)
		spawn(NpcPot, 20)
		spawn(NpcFlower, 20)
		spawn(NpcCloud, 10)
		step = 11
	}

	/*
		spawnerCooldown = time.Now().Add(3 * time.Second)
		if len(npcs) > spawnMax {
			return
		}
		spawnCount := spawnMax + 1 - len(npcs)
		if spawnCount > 3 {
			spawnerCooldown = time.Now().Add(500 * time.Millisecond)
			spawnCount = 3
		}

		for i := 0; i < spawnCount; i++ {

			theta := rand.Float64() * 6

			distance := minDistance + (rand.Float64() * (maxDistance - minDistance))

			New(rand.Intn(6-1)+1, global.Player.X()+math.Sin(theta)*distance, global.Player.Y()+math.Cos(theta)*distance, global.Player)
		}
		log.Debug().Msgf("spawned %d", spawnCount)
	*/
}

func spawn(npcType int, count int) {
	for i := 0; i < count; i++ {
		theta := rand.Float64() * 6
		distance := minDistance + (rand.Float64() * (maxDistance - minDistance))
		New(npcType, global.Player.X()+math.Sin(theta)*distance, global.Player.Y()+math.Cos(theta)*distance, global.Player)
	}
}

func Knockback(distance float64) {
	for _, n := range npcs {
		fmt.Println(global.Distance(global.Player.X(), global.Player.Y(), n.x, n.y))
		if global.Distance(global.Player.X(), global.Player.Y(), n.x, n.y) >= 100 {
			continue
		}
		if n.x > global.Player.X() {
			n.x += distance
		} else if n.x <= global.Player.X() {
			n.x -= distance
		}
		if n.y > global.Player.Y() {
			n.y += distance
		} else if n.y <= global.Player.Y() {
			n.y -= distance
		}
	}
}
