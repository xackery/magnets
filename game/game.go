package game

import (
	"context"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/xackery/magnets/bullet"
	"github.com/xackery/magnets/camera"
	"github.com/xackery/magnets/collision"
	"github.com/xackery/magnets/font"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/input"
	"github.com/xackery/magnets/library"
	"github.com/xackery/magnets/npc"
	"github.com/xackery/magnets/player"
	"github.com/xackery/magnets/weapon"
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
	ctx              context.Context
	cancel           context.CancelFunc
	resolutionChange time.Time
	frame            int64
	nextFrame        time.Time
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

	npc.New(npc.NpcBat, 20, 0)
	npc.New(npc.NpcCloud, 40, 0)
	npc.New(npc.NpcFlower, 60, 0)
	npc.New(npc.NpcAseprite, 80, 0)
	npc.New(npc.NpcPot, 100, 0)
	npc.New(npc.NpcKnight, 120, 0)

	p, err := player.New("player", "default")
	if err != nil {
		return nil, fmt.Errorf("player new player: %w", err)
	}

	weapons := []int{
		//weapon.WeaponArrow,
		//weapon.WeaponBoomerang,
		//weapon.WeaponShuriken,
		//weapon.WeaponSpear,
		//weapon.WeaponSword,
	}
	for _, w := range weapons {
		wp, _ := weapon.New(w)
		p.WeaponAdd(wp)
	}

	g.Layout(global.ScreenWidth(), global.ScreenHeight())
	Instance = g
	return g, nil
}

// Draw is called for render update
func (g *Game) Draw(screen *ebiten.Image) {
	collision.Image = ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	screen.Fill(color.RGBA64{R: 0x5050, G: 0x5050, B: 0xcfcf, A: 0xFFFF})
	//update game
	//x, y := ebiten.CursorPosition()

	npc.Draw(collision.Image)
	bullet.Draw(screen)
	player.Draw(collision.Image)

	screen.DrawImage(collision.Image, nil)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f, Position: %0.2f, %0.2f, Camera: %0.2f, %0.2f", ebiten.CurrentTPS(), player.X(), player.Y(), camera.X, camera.Y), 0, global.ScreenHeight()-14)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(x, y int) (int, int) {
	global.ScreenOnLayoutChange(x, y)
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
		g.frame = 0
	}

	ebiten.CurrentTPS()

	input.Update()
	bullet.Update()
	player.Update()
	if !player.IsDead() {
		npc.Update(player.X(), player.Y())
	}

	if g.frame%10 == 0 {
		bullet.HitUpdate()
		player.HitUpdate()
		npc.HitUpdate()
	}
	return nil
}
