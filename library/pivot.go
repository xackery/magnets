package library

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/xackery/aseprite"
)

// Pivot represents attachment or pivot points
type Pivot struct {
	Cells []*PivotCell
}

// PivotCell represents each unique animation frame of a pivot
type PivotCell struct {
	Pivots map[int]*PivotData
}

// PivotData represents the inner pivot details
type PivotData struct {
	PositionX           int
	PositionY           int
	IsBehind            bool
	IsHorizontalFlipped bool
	IsVerticalFlipped   bool
}

// PivotInfo returns specified pivot details, used by sprites
func PivotInfo(spriteName string, cellIndex int, pivotIndex int) (*PivotData, error) {
	pivot := pivots[spriteName]
	if pivot == nil {
		return nil, fmt.Errorf("%s not found", spriteName)
	}
	if len(pivot.Cells) < cellIndex {
		return nil, fmt.Errorf("cell %d not found", cellIndex)
	}

	pivotData := pivot.Cells[cellIndex].Pivots[pivotIndex]
	if pivotData == nil {
		return nil, fmt.Errorf("pivotIndex %d not found", pivotIndex)
	}

	return pivotData, nil
}

func assessPivots(assetName string, l *aseprite.Layer, c *aseprite.Cell, cellIndex int) {
	if strings.ToLower(l.Name) != "pivot" {
		return
	}

	ext := filepath.Ext(assetName)
	baseName := assetName[0 : len(assetName)-len(ext)]

	pivot := pivots[baseName]
	if pivot == nil {
		pivot = &Pivot{}
		pivots[baseName] = pivot
		log.Debug().Msgf("added pivots for %s", baseName)
	}

	//pad cells to ensure any empty cells are accounted for
	for i := 0; i <= cellIndex; i++ {
		pivot.Cells = append(pivot.Cells, &PivotCell{
			Pivots: make(map[int]*PivotData),
		})
	}

	pivotCell := pivot.Cells[cellIndex]

	for x := 0; x < c.Image.Rect.Max.X; x++ {
		for y := 0; y < c.Image.Rect.Max.Y; y++ {
			r, g, b, a := c.Image.At(x, y).RGBA()
			for key, val := range []uint32{r, g, b} {
				if val == 0 {
					continue
				}

				pivotData := &PivotData{
					PositionX: x + int(c.PositionX),
					PositionY: y + int(c.PositionY),
				}

				if val == 0xFEFE {
					pivotData.IsHorizontalFlipped = true
				}
				if val == 0xFDFD {
					pivotData.IsVerticalFlipped = true
				}
				if val == 0xFCFC {
					pivotData.IsHorizontalFlipped = true
					pivotData.IsVerticalFlipped = true
				}
				if a != 0xFFFF {
					pivotData.IsBehind = true
				}

				pivotCell.Pivots[key] = pivotData
			}
		}
	}
}
