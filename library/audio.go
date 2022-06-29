package library

import (
	"fmt"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/xackery/magnets/global"
)

var (
	audios = make(map[string]*AudioData)
)

// AudioData represents some form of audio playback
type AudioData struct {
	Player *audio.Player
}

func AudioPlay(name string) error {
	if runtime.GOOS == "js" {
		return nil
	}
	a, err := Audio(name)
	if err != nil {
		return err
	}
	a.Player.Rewind()
	a.Player.Play()
	return nil
}

// Audio returns an audio represented by name
func Audio(name string) (*AudioData, error) {
	a := audios[name]
	if a == nil {
		return nil, fmt.Errorf("audio not found")
	}
	return a, nil
}

func onVolumeChange() {
	for _, a := range audios {
		a.Player.SetVolume(global.Volume())
	}
}
