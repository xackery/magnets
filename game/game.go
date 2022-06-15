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
	"github.com/xackery/magnets/font"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
	"github.com/xackery/magnets/npc"
	"github.com/xackery/magnets/player"
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
	bulletTimer      *time.Ticker
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
	//window resizable
	ebiten.SetWindowResizable(true)

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

	_, err = npc.New("heart", "default", 0, global.AnchorTopLeft, 250, 30)
	if err != nil {
		return nil, fmt.Errorf("npc new heart: %w", err)
	}
	_, err = player.New("player", "default", 0, global.AnchorTopLeft, 120, 30)
	if err != nil {
		return nil, fmt.Errorf("player new player: %w", err)
	}

	bullet.New(bullet.BulletPlain, 0, global.AnchorTopLeft, 150, 60)
	g.bulletTimer = time.NewTicker(250 * time.Millisecond)

	Instance = g
	return g, nil
}

// Draw is called for render update
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA64{R: 0x5050, G: 0x5050, B: 0xcfcf, A: 0xFFFF})
	//update game
	x, y := ebiten.CursorPosition()

	/*select {
	case <-g.bulletTimer.C:
		_, err := bullet.New(bullet.BulletPlain, 0, global.AnchorTopLeft, 150, 30)
		if err != nil {
			fmt.Println("bullet creation failed", err)
		}
	default:
	}*/
	npc.Draw(screen)
	bullet.Draw(screen)
	player.Draw(screen)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f, Position: %d, %d", ebiten.CurrentTPS(), x, y), 0, global.ScreenHeight()-14)

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

	bullet.Update()
	return nil
}
