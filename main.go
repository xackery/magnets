package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/xackery/magnets/game"
	"github.com/xackery/magnets/global"
	"github.com/xackery/magnets/library"
)

var (
	// Version of the build
	Version string
	// BuildDate is when the last build happened
	BuildDate string
	// WebAssetPath represents the absolute URL of where web assets are for js builds
	WebAssetPath string
	// ScreenWidth is the width of the resolution
	ScreenWidth = 640 //640 //1280
	// ScreenHeight is the height of the resolution
	ScreenHeight = 480 //480 //720
)

func main() {
	var err error
	stage := "alpha"
	if runtime.GOOS != "windows" {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	tag := fmt.Sprintf("%s %s", Version, BuildDate)
	if Version == "" {
		tag = time.Now().Format("2006-01-02")
	}
	ebiten.SetWindowTitle(fmt.Sprintf("magnets %s %s", stage, tag))
	log.Info().Msgf("starting magnets %s %s %s", stage, Version, BuildDate)
	/*if runtime.GOOS == "js" {
		query := js.Global().Get("window").Get("location").Get("search")
		param := fmt.Sprintf("%s", query)
		if strings.HasPrefix(param, "?") {
			param = param[1:]
		}
		log.Debug().Msgf("query: %s", param)
		if len(param) > 0 {
			values, err := url.ParseQuery(param)
			if err != nil {
				log.Debug().Err(err).Msg("failed to parse query")
			} else {
				log.Debug().Msgf("room: %s", values.Get("room"))
			}
		}
	}*/

	if runtime.GOOS != "darwin" {
		global.ScreenScaleChange(4, 4)
	}

	//clipboard.WriteAll(runtime.GOOS)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := time.Now()
	if len(WebAssetPath) > 0 {
		library.SetWebAssetPath(WebAssetPath)
	}
	global.ScreenOnLayoutChange(ScreenWidth, ScreenHeight, true)
	global.ScreenScaleChange(ebiten.DeviceScaleFactor(), ebiten.DeviceScaleFactor())
	g, err := game.New(ctx, WebAssetPath)
	if err != nil {
		log.Error().Err(err).Msg("game")
		os.Exit(1)
	}
	log.Info().Msgf("game loaded in %0.2f seconds", time.Since(start).Seconds())
	err = ebiten.RunGame(g)
	if err != nil {
		log.Error().Err(err).Msg("runGame")
		if runtime.GOOS == "windows" {
			option := ""
			log.Info().Msg("press a key then enter to exit.")
			fmt.Scan(&option)
		}
		os.Exit(1)
	}

	log.Info().Msgf("exited safely after %0.2f seconds", time.Since(start).Seconds())
}
