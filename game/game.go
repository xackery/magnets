package game

import (
	"context"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/rs/zerolog/log"
	"github.com/xackery/magnets/bullet"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/collision"
	"github.com/xackery/magnets/font"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/input"
	"github.com/xackery/magnets/item"
	"github.com/xackery/magnets/library"
	"github.com/xackery/magnets/npc"
	"github.com/xackery/magnets/player"
	"github.com/xackery/magnets/ui/bar"
	"github.com/xackery/magnets/ui/equipment"
	"github.com/xackery/magnets/ui/gameover"
	"github.com/xackery/magnets/ui/level"
	"github.com/xackery/magnets/ui/life"
	"github.com/xackery/magnets/world"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	Instance *Game
)

const (
	StageBattle = iota
	StageReward
	StageOverworld
)

// Game implements the ebiten Game interface
type Game struct {
	ctx               context.Context
	cancel            context.CancelFunc
	resolutionChange  time.Time
	frame             int64
	nextFrame         time.Time
	resetGameCooldown time.Time
}

// New creates a new game instance
func New(ctx context.Context, host string) (*Game, error) {
	var err error
	ctx, cancel := context.WithCancel(ctx)
	g := &Game{
		ctx:    ctx,
		cancel: cancel,
	}
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(global.ScreenWidth(), global.ScreenHeight())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	//ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)

	g.resolutionChange = time.Now().Add(1 * time.Second)
	err = font.Load()
	if err != nil {
		return nil, fmt.Errorf("font.Load: %w", err)
	}

	err = library.LoadFontTTF("goregular", goregular.TTF, nil, 'M')
	if err != nil {
		return nil, fmt.Errorf("loadfont goregular: %w", err)
	}
	//true fullscreen
	//ebiten.SetFullscreen(true)

	//fullscreen borderless
	//ebiten.SetWindowDecorated(false)
	//ebiten.SetWindowSize(ebiten.MonitorSize())

	for _, file := range files {
		err = library.Load(file)
		if err != nil {
			if file == "state.yml" {
				continue
			}
			return nil, fmt.Errorf("library.Load: %w", err)
		}
	}

	g.clear()
	err = g.start()
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	Instance = g
	return g, nil
}

func (g *Game) clear() {
	global.SetIsPaused(false)
	global.Kill = 0
	npc.Clear()
	player.Clear()
	world.Clear()
	bullet.Clear()
	item.Clear()
	equipment.Clear()
	bar.Clear()
	gameover.Clear()
	global.SetIsGameOver(false)
	global.SetIsLevelUp(false)
}

func (g *Game) start() error {
	global.Countdown = 600
	_, err := player.New("player", "base")
	if err != nil {
		return fmt.Errorf("player new hero: %w", err)
	}
	/*npc.New(npc.NpcBat, 20, 0, p)
	npc.New(npc.NpcCloud, 40, 0, p)
	npc.New(npc.NpcFlower, 60, 0, p)
	npc.New(npc.NpcAseprite, 80, 0, p)
	npc.New(npc.NpcPot, 100, 0, p)
	npc.New(npc.NpcKnight, 120, 0, p)*/

	/*_, err = world.New(world.WorldGrass)
	if err != nil {
		return fmt.Errorf("world.New: %w", err)
	}*/

	global.ScreenOnLayoutChange(global.ScreenWidth(), global.ScreenHeight(), true)
	return nil
}

// Draw is called for render update
func (g *Game) Draw(screen *ebiten.Image) {
	collision.Image = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	screen.Fill(color.RGBA64{R: 6565, G: 0x9090, B: 0x3d3d, A: 0xFFFF})
	//update game
	//x, y := ebiten.CursorPosition()

	//world.Draw(screen)
	npc.Draw(collision.Image)
	item.Draw(screen)
	bullet.Draw(screen)
	player.Draw(collision.Image)

	screen.DrawImage(collision.Image, nil)
	bar.Draw(screen)
	equipment.Draw(screen)
	life.Draw(screen)
	level.Draw(screen)
	gameover.Draw(screen)

	g.countdownDraw(screen)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f, Position: %0.2f, %0.2f, Camera: %0.2f, %0.2f", ebiten.CurrentTPS(), player.X(), player.Y(), camera.X, camera.Y), 0, global.ScreenHeight()-14)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(x, y int) (int, int) {
	global.ScreenOnLayoutChange(x, y, false)
	return x, y
}

// Update is for each update game
func (g *Game) Update() error {
	if time.Now().Before(g.resolutionChange) {
		g.resolutionChange = time.Now().Add(100 * time.Hour)
		//ebiten.SetWindowSize(640, 480)
	}
	g.frame++
	if time.Now().After(g.nextFrame) {
		g.nextFrame = time.Now().Add(1 * time.Second)
		if !global.IsPaused() {
			global.Countdown--
		}
		g.frame = 0
	}

	ebiten.CurrentTPS()

	if ebiten.IsKeyPressed(ebiten.KeyP) && time.Now().After(g.resetGameCooldown) {
		g.resetGameCooldown = time.Now().Add(1 * time.Second)
		global.SetIsPaused(!global.IsPaused())
	}
	input.Update()

	if ebiten.IsKeyPressed(ebiten.KeyR) && time.Now().After(g.resetGameCooldown) {
		g.resetGameCooldown = time.Now().Add(3 * time.Second)
		g.clear()
		err := g.start()
		if err != nil {
			log.Debug().Err(err).Msgf("start")
		}

		return nil
	}

	player.Update()
	if global.IsPaused() {
		return nil
	}
	bullet.Update()
	item.Update()
	if !player.IsDead() {
		npc.Update()
	}

	if g.frame%2 == 0 {
		bullet.HitUpdate()
		player.HitUpdate()
		npc.HitUpdate()
	}
	return nil
}

func (g *Game) countdownDraw(screen *ebiten.Image) {
	width := global.ScreenWidth()
	height := 64
	msg := fmt.Sprintf("Survive for %d seconds!", global.Countdown)

	bounds := text.BoundString(font.NormalFont(), msg)
	x, y := int(width/2)-bounds.Min.X-bounds.Dx()/2, int(height/2)-bounds.Min.Y-bounds.Dy()/2

	text.Draw(screen, msg, font.NormalFont(), x-1, y-1, color.Black)
	text.Draw(screen, msg, font.NormalFont(), x, y, color.White)

	x += 350
	y += 50

	text.Draw(screen, fmt.Sprintf("%d kills", global.Kill), font.NormalFont(), x-1, y-1, color.Black)
	text.Draw(screen, fmt.Sprintf("%d kills", global.Kill), font.NormalFont(), x, y, color.White)

}
