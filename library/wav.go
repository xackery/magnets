package library

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

func loadWav(assetName string) (*AudioData, error) {
	ext := filepath.Ext(assetName)
	baseName := assetName[0 : len(assetName)-len(ext)]
	audioName := baseName + ".wav"

	r, err := ReadFile(audioName)
	if err != nil {
		return nil, fmt.Errorf("readFile %s: %w", audioName, err)
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("readall: %w", err)
	}

	s, err := wav.Decode(audioContext, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	p, err := audio.NewPlayer(audioContext, s)
	if err != nil {
		return nil, fmt.Errorf("newplayer: %w", err)
	}
	a := &AudioData{
		Player: p,
	}
	return a, nil
}
