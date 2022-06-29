package library

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/rs/zerolog/log"
	"github.com/xackery/aseprite"
	"github.com/xackery/magnets/art"
	"github.com/xackery/magnets/global"
)

var (
	images        = make(map[string]*ebiten.Image)
	markedSprites []string
	pivots        = make(map[string]*Pivot)
	fonts         = make(map[string]*FontData)
	webAssetPath  string
	mu            sync.RWMutex
	audioContext  *audio.Context
)

func Init() {
	global.SubscribeOnVolumeChange(onVolumeChange)
}

// Flush will flush all sprites not in markedSprites.
// It should be called after scene finishes loading.
func Flush() {
	deleteCount := 0
	for name := range sprites {
		isMarked := false
		for _, markName := range markedSprites {
			if markName != name {
				continue
			}
			isMarked = true
			break
		}
		if !isMarked {
			delete(sprites, name)
			deleteCount++
		}
	}

	log.Debug().Msgf("removed %d assets from memory, %d retained", deleteCount, len(sprites))
	markedSprites = []string{}
}

// Load should be called during initialization of a scene
// Each loaded asset is put in markedSprites.
// After finishing loading a scene, use Flush() to purge unused sprites
func Load(assetName string) error {
	if audioContext == nil {
		audioContext = audio.NewContext(22050)
	}
	assetName = strings.ToLower(assetName)

	ext := filepath.Ext(assetName)
	baseName := assetName[0 : len(assetName)-len(ext)]
	baseName = filepath.Base(baseName)
	switch ext {
	case ".aseprite":
		sprite, err := loadAseprite(assetName)
		if err != nil {
			return fmt.Errorf("loadAseprite: %w", err)
		}
		sprites[baseName] = sprite
	case ".png":
		image, err := loadPng(assetName)
		if err != nil {
			return fmt.Errorf("loadPng: %w", err)
		}
		images[baseName] = image
	case ".ogg":
		a, err := loadOgg(assetName)
		if err != nil {
			return fmt.Errorf("loadOgg: %w", err)
		}
		a.Player.SetVolume(global.Volume())
		audios[baseName] = a
	case ".wav":
		if runtime.GOOS == "js" {
			return nil
		}
		a, err := loadWav(assetName)
		if err != nil {
			return fmt.Errorf("loadWav: %w", err)
		}
		a.Player.SetVolume(global.Volume())
		audios[baseName] = a
	default:
		return fmt.Errorf("unsupported extension: %s", ext)
	}

	markedSprites = append(markedSprites, baseName)
	return nil
}

// Layer returns requested layer if available
func Layer(spriteName string, layerName string) (*aseprite.Layer, error) {
	spriteName = strings.ToLower(spriteName)
	sprite, ok := sprites[spriteName]
	if !ok {
		return nil, fmt.Errorf("sprite %s not found", spriteName)
	}
	layerName = strings.ToLower(layerName)
	layer, ok := sprite.Layers[layerName]
	if !ok {
		return nil, fmt.Errorf("layer %s not found", layerName)
	}
	return layer, nil
}

func Image(imageName string) (*ebiten.Image, error) {
	imageName = strings.ToLower(imageName)
	image, ok := images[imageName]
	if !ok {
		return nil, fmt.Errorf("image %s not found", imageName)
	}
	return image, nil
}

// Tag returns requested tag if available
func Tag(spriteName string, tagName string) (*aseprite.Tag, error) {
	spriteName = strings.ToLower(spriteName)
	sprite, ok := sprites[spriteName]
	if !ok {
		return nil, fmt.Errorf("sprite %s not found", spriteName)
	}
	tagName = strings.ToLower(tagName)
	for _, tag := range sprite.Tags {
		if strings.ToLower(tag.Name) == tagName {
			return tag, nil
		}
	}
	return nil, fmt.Errorf("tag %s not found", tagName)
}

// SetWebAssetPath sets the web asset URL path
func SetWebAssetPath(path string) {
	mu.Lock()
	webAssetPath = path
	mu.Unlock()
}

func ReadFile(assetName string) (assetReader, error) {
	//if runtime.GOOS != "js" {
	data, err := os.ReadFile(fmt.Sprintf("art/%s", assetName))
	if err != nil {
		data, err = art.Art.ReadFile(assetName)
		if err != nil {
			return nil, fmt.Errorf("readFile %s: %w", assetName, err)
		}
	}
	r := bytes.NewReader(data)
	return r, nil
	//}
	/*
		resp, err := http.Get(fmt.Sprintf("%s/art/%s", webAssetPath, assetName))
		if err != nil {
			return nil, fmt.Errorf("art get: %w", err)
		}
		defer resp.Body.Close()
		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, resp.Body)
		if err != nil {
			return nil, fmt.Errorf("copy asset: %w", err)
		}
		r = bytes.NewReader(buf.Bytes())
		return r, nil
	*/
}
